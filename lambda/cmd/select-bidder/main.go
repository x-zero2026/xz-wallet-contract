package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/x-zero/xz-wallet/pkg/auth"
	"github.com/x-zero/xz-wallet/pkg/blockchain"
	"github.com/x-zero/xz-wallet/pkg/db"
	"github.com/x-zero/xz-wallet/pkg/response"
)

type SelectBidderRequest struct {
	BidderDID string `json:"bidder_did"`
}

type SelectBidderResponse struct {
	TaskID      string `json:"task_id"`
	ExecutorDID string `json:"executor_did"`
	TxHash      string `json:"tx_hash"`
	Status      string `json:"status"`
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

	var req SelectBidderRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.Error(400, "Invalid request body")
	}

	if req.BidderDID == "" {
		return response.Error(400, "Missing bidder_did")
	}

	if err := db.InitDB(); err != nil {
		return response.Error(500, fmt.Sprintf("Database error: %v", err))
	}

	client, err := blockchain.InitClient()
	if err != nil {
		return response.Error(500, fmt.Sprintf("Blockchain error: %v", err))
	}

	pool := db.GetPool()

	// Get task
	var task struct {
		ContractTaskID int64
		CreatorDID     string
		Status         string
	}
	err = pool.QueryRow(ctx, `
		SELECT contract_task_id, creator_did, status FROM tasks WHERE task_id = $1
	`, taskID).Scan(&task.ContractTaskID, &task.CreatorDID, &task.Status)
	if err != nil {
		return response.Error(404, "Task not found")
	}

	// Verify user is creator
	if task.CreatorDID != claims.DID {
		return response.Error(403, "Only creator can select bidder")
	}

	// Verify task status
	if task.Status != "bidding" {
		return response.Error(400, "Task is not in bidding status")
	}

	// Verify bid exists
	var bidID string
	err = pool.QueryRow(ctx, `
		SELECT bid_id FROM task_bids WHERE task_id = $1 AND bidder_did = $2 AND status = 'pending'
	`, taskID, req.BidderDID).Scan(&bidID)
	if err != nil {
		return response.Error(404, "Bid not found")
	}

	// Get bidder's eth_address
	var executorEthAddress string
	err = pool.QueryRow(ctx, "SELECT eth_address FROM users WHERE did = $1", req.BidderDID).Scan(&executorEthAddress)
	if err != nil {
		return response.Error(404, "Bidder not found")
	}

	// Set executor on blockchain
	txHash, err := client.SetExecutor(uint64(task.ContractTaskID), executorEthAddress)
	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to set executor on blockchain: %v", err))
	}

	// Update database
	tx, err := pool.Begin(ctx)
	if err != nil {
		return response.Error(500, "Failed to start transaction")
	}
	defer tx.Rollback(ctx)

	// Update task
	_, err = tx.Exec(ctx, `
		UPDATE tasks SET executor_did = $1, status = 'accepted', updated_at = CURRENT_TIMESTAMP
		WHERE task_id = $2
	`, req.BidderDID, taskID)
	if err != nil {
		return response.Error(500, "Failed to update task")
	}

	// Update selected bid
	_, err = tx.Exec(ctx, `
		UPDATE task_bids SET status = 'accepted', updated_at = CURRENT_TIMESTAMP
		WHERE bid_id = $1
	`, bidID)
	if err != nil {
		return response.Error(500, "Failed to update bid")
	}

	// Reject other bids
	_, err = tx.Exec(ctx, `
		UPDATE task_bids SET status = 'rejected', updated_at = CURRENT_TIMESTAMP
		WHERE task_id = $1 AND bidder_did != $2 AND status = 'pending'
	`, taskID, req.BidderDID)
	if err != nil {
		return response.Error(500, "Failed to reject other bids")
	}

	if err := tx.Commit(ctx); err != nil {
		return response.Error(500, "Failed to commit transaction")
	}

	return response.Success(SelectBidderResponse{
		TaskID:      taskID,
		ExecutorDID: req.BidderDID,
		TxHash:      txHash,
		Status:      "accepted",
	})
}

func main() {
	lambda.Start(handler)
}
