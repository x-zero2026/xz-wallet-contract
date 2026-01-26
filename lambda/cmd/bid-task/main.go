package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/x-zero/xz-wallet/pkg/auth"
	"github.com/x-zero/xz-wallet/pkg/db"
	"github.com/x-zero/xz-wallet/pkg/response"
)

type BidTaskRequest struct {
	Message string `json:"message"`
}

type BidTaskResponse struct {
	BidID  string `json:"bid_id"`
	Status string `json:"status"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type,Authorization",
				"Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
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

	var req BidTaskRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.Error(400, "Invalid request body")
	}

	if err := db.InitDB(); err != nil {
		return response.Error(500, fmt.Sprintf("Database error: %v", err))
	}

	pool := db.GetPool()

	// Get user's credit score
	var creditScore int
	err = pool.QueryRow(ctx, "SELECT credit_score FROM users WHERE did = $1", claims.DID).Scan(&creditScore)
	if err != nil {
		return response.Error(404, "User not found")
	}

	// Check credit score
	if creditScore < 0 {
		return response.Error(403, "Insufficient credit score to bid")
	}

	// Check task exists and is biddable
	var taskStatus, creatorDID string
	err = pool.QueryRow(ctx, "SELECT status, creator_did FROM tasks WHERE task_id = $1", taskID).Scan(&taskStatus, &creatorDID)
	if err != nil {
		return response.Error(404, "Task not found")
	}

	if creatorDID == claims.DID {
		return response.Error(403, "Cannot bid on your own task")
	}

	if taskStatus != "pending" && taskStatus != "bidding" {
		return response.Error(400, "Task is not accepting bids")
	}

	// Insert bid
	var bidID string
	err = pool.QueryRow(ctx, `
		INSERT INTO task_bids (task_id, bidder_did, bid_message, credit_score_snapshot, status)
		VALUES ($1, $2, $3, $4, 'pending')
		ON CONFLICT (task_id, bidder_did) DO UPDATE
		SET bid_message = $3, updated_at = CURRENT_TIMESTAMP
		RETURNING bid_id
	`, taskID, claims.DID, req.Message, creditScore).Scan(&bidID)
	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to create bid: %v", err))
	}

	// Update task status to bidding if it was pending
	if taskStatus == "pending" {
		_, err = pool.Exec(ctx, "UPDATE tasks SET status = 'bidding' WHERE task_id = $1", taskID)
		if err != nil {
			return response.Error(500, "Failed to update task status")
		}
	}

	return response.Success(BidTaskResponse{
		BidID:  bidID,
		Status: "pending",
	})
}

func main() {
	lambda.Start(handler)
}
