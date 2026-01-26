package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/x-zero/xz-wallet/pkg/auth"
	"github.com/x-zero/xz-wallet/pkg/db"
	"github.com/x-zero/xz-wallet/pkg/models"
	"github.com/x-zero/xz-wallet/pkg/response"
)

type SubmitWorkRequest struct {
	SubmissionType string   `json:"submission_type"` // "design", "implementation", "final"
	Content        string   `json:"content"`
	FileURLs       []string `json:"file_urls,omitempty"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type,Authorization",
				"Access-Control-Allow-Methods": "POST,OPTIONS",
			},
		}, nil
	}

	// Validate JWT
	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		authHeader = request.Headers["authorization"]
	}
	claims, err := auth.ValidateToken(authHeader)
	if err != nil {
		return response.Error(401, fmt.Sprintf("Invalid token: %v", err))
	}

	taskID := request.PathParameters["id"]
	if taskID == "" {
		return response.Error(400, "Missing task ID")
	}

	var req SubmitWorkRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.Error(400, "Invalid request body")
	}

	if req.SubmissionType == "" || req.Content == "" {
		return response.Error(400, "Missing submission_type or content")
	}

	if err := db.InitDB(); err != nil {
		return response.Error(500, fmt.Sprintf("Database error: %v", err))
	}

	pool := db.GetPool()

	// Get task and verify executor
	var task struct {
		ExecutorDID *string
		Status      string
	}
	err = pool.QueryRow(ctx, `
		SELECT executor_did, status FROM tasks WHERE task_id = $1
	`, taskID).Scan(&task.ExecutorDID, &task.Status)
	if err != nil {
		return response.Error(404, "Task not found")
	}

	// Verify user is the executor
	if task.ExecutorDID == nil || *task.ExecutorDID != claims.DID {
		return response.Error(403, "Only task executor can submit work")
	}

	// Validate status and determine new status
	var newStatus string
	switch req.SubmissionType {
	case "design":
		if task.Status != models.TaskStatusAccepted {
			return response.Error(400, "Can only submit design when task is accepted")
		}
		newStatus = models.TaskStatusDesignSubmitted
	case "implementation":
		if task.Status != models.TaskStatusDesignApproved {
			return response.Error(400, "Can only submit implementation after design is approved")
		}
		newStatus = models.TaskStatusImplementationSubmitted
	case "final":
		if task.Status != models.TaskStatusImplementationApproved {
			return response.Error(400, "Can only submit final work after implementation is approved")
		}
		newStatus = models.TaskStatusFinalSubmitted
	default:
		return response.Error(400, "Invalid submission_type. Must be: design, implementation, or final")
	}

	// Start transaction
	tx, err := pool.Begin(ctx)
	if err != nil {
		return response.Error(500, "Failed to start transaction")
	}
	defer tx.Rollback(ctx)

	// Insert submission record
	var submissionID string
	err = tx.QueryRow(ctx, `
		INSERT INTO task_submissions (task_id, submission_type, content, file_urls, status)
		VALUES ($1, $2, $3, $4, 'pending')
		RETURNING submission_id
	`, taskID, req.SubmissionType, req.Content, req.FileURLs).Scan(&submissionID)
	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to create submission: %v", err))
	}

	// Update task status
	_, err = tx.Exec(ctx, `
		UPDATE tasks 
		SET status = $1, updated_at = NOW()
		WHERE task_id = $2
	`, newStatus, taskID)
	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to update task status: %v", err))
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return response.Error(500, "Failed to commit transaction")
	}

	return response.Success(map[string]interface{}{
		"message":       "Work submitted successfully",
		"submission_id": submissionID,
		"new_status":    newStatus,
	})
}

func main() {
	lambda.Start(handler)
}
