package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/x-zero/xz-wallet/pkg/auth"
	"github.com/x-zero/xz-wallet/pkg/blockchain"
	"github.com/x-zero/xz-wallet/pkg/db"
	"github.com/x-zero/xz-wallet/pkg/response"
)

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

	if err := db.InitDB(); err != nil {
		return response.Error(500, fmt.Sprintf("Database error: %v", err))
	}

	pool := db.GetPool()

	// Get task and verify ownership
	var task struct {
		ContractTaskID int64
		CreatorDID     string
		ExecutorDID    *string
		Status         string
		PaidAmount     string
	}
	err = pool.QueryRow(ctx, `
		SELECT contract_task_id, creator_did, executor_did, status, paid_amount 
		FROM tasks WHERE task_id = $1
	`, taskID).Scan(&task.ContractTaskID, &task.CreatorDID, &task.ExecutorDID, &task.Status, &task.PaidAmount)
	if err != nil {
		return response.Error(404, "Task not found")
	}

	// Verify user is the creator or executor
	isCreator := task.CreatorDID == claims.DID
	isExecutor := task.ExecutorDID != nil && *task.ExecutorDID == claims.DID
	
	if !isCreator && !isExecutor {
		return response.Error(403, "Only task creator or executor can cancel the task")
	}

	// Only allow cancellation if task is not completed or already cancelled
	if task.Status == "completed" || task.Status == "cancelled" {
		return response.Error(400, fmt.Sprintf("Cannot cancel task in status: %s", task.Status))
	}

	// Initialize blockchain client
	client, err := blockchain.InitClient()
	if err != nil {
		return response.Error(500, fmt.Sprintf("Blockchain error: %v", err))
	}

	// Calculate executor amount for cancellation
	// Note: executorAmount is the ADDITIONAL amount to pay executor during cancellation
	// Since executor has already been paid (task.PaidAmount), we don't pay them again
	// They keep what they've already received
	executorAmount := big.NewInt(0)

	// Cancel task on blockchain
	txHash, err := client.CancelTask(uint64(task.ContractTaskID), executorAmount)
	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to cancel task on blockchain: %v", err))
	}

	// Update task status to cancelled
	_, err = pool.Exec(ctx, `
		UPDATE tasks 
		SET status = 'cancelled',
		    cancelled_at = NOW(),
		    updated_at = NOW()
		WHERE task_id = $1
	`, taskID)
	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to cancel task: %v", err))
	}

	// Apply credit score penalty if executor quits mid-task
	if isExecutor && task.ExecutorDID != nil {
		var creditPenalty int
		
		// Determine penalty based on task status
		switch task.Status {
		case "design_approved", "implementation_submitted":
			// After design payment, before implementation payment: -100
			creditPenalty = 100
		case "implementation_approved", "final_submitted":
			// After implementation payment, before final payment: -200
			creditPenalty = 200
		default:
			// Other stages: no penalty
			creditPenalty = 0
		}
		
		if creditPenalty > 0 {
			_, err = pool.Exec(ctx, `
				UPDATE users 
				SET credit_score = credit_score - $1,
				    tasks_cancelled = tasks_cancelled + 1
				WHERE did = $2
			`, creditPenalty, *task.ExecutorDID)
			if err != nil {
				return response.Error(500, fmt.Sprintf("Failed to update executor credit: %v", err))
			}
		}
	}

	return response.Success(map[string]interface{}{
		"message": "Task cancelled successfully",
		"task_id": taskID,
		"tx_hash": txHash,
	})
}

func main() {
	lambda.Start(handler)
}
