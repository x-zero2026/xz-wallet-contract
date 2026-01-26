package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/x-zero/xz-wallet/pkg/auth"
	"github.com/x-zero/xz-wallet/pkg/blockchain"
	"github.com/x-zero/xz-wallet/pkg/db"
	"github.com/x-zero/xz-wallet/pkg/models"
	"github.com/x-zero/xz-wallet/pkg/response"
)

type ApproveWorkRequest struct {
	Milestone string `json:"milestone"` // "design", "implementation", "final"
	Approve   bool   `json:"approve"`   // true for approve, false for reject
	RejectionReason string `json:"rejection_reason,omitempty"` // optional reason for rejection
}

type ApproveWorkResponse struct {
	Status  string        `json:"status"`
	Payment PaymentDetail `json:"payment"`
}

type PaymentDetail struct {
	Amount string `json:"amount"`
	TxHash string `json:"tx_hash"`
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

	var req ApproveWorkRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.Error(400, "Invalid request body")
	}

	// Initialize
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
		RewardAmount   string
		PaidAmount     string
	}
	err = pool.QueryRow(ctx, `
		SELECT contract_task_id, creator_did, status, reward_amount, paid_amount
		FROM tasks WHERE task_id = $1
	`, taskID).Scan(&task.ContractTaskID, &task.CreatorDID, &task.Status, &task.RewardAmount, &task.PaidAmount)
	if err != nil {
		return response.Error(404, "Task not found")
	}

	// Verify user is creator
	if task.CreatorDID != claims.DID {
		return response.Error(403, "Only creator can approve work")
	}

	// Determine current and target status based on milestone
	var currentStatus, approvedStatus, rejectedStatus string
	var paymentPercent int
	
	switch req.Milestone {
	case "design":
		currentStatus = models.TaskStatusDesignSubmitted
		approvedStatus = models.TaskStatusDesignApproved
		rejectedStatus = models.TaskStatusAccepted // Reject back to accepted
		paymentPercent = models.MilestoneDesign // 30%
	case "implementation":
		currentStatus = models.TaskStatusImplementationSubmitted
		approvedStatus = models.TaskStatusImplementationApproved
		rejectedStatus = models.TaskStatusDesignApproved // Reject back to design approved
		paymentPercent = models.MilestoneImplementation // 50%
	case "final":
		currentStatus = models.TaskStatusFinalSubmitted
		approvedStatus = models.TaskStatusCompleted
		rejectedStatus = models.TaskStatusImplementationApproved // Reject back to implementation approved
		paymentPercent = models.MilestoneFinal // 20%
	default:
		return response.Error(400, "Invalid milestone")
	}

	// Verify task status
	if task.Status != currentStatus {
		return response.Error(400, fmt.Sprintf("Invalid task status for %s approval", req.Milestone))
	}

	// Handle rejection
	if !req.Approve {
		// Update submission status to rejected
		_, err = pool.Exec(ctx, `
			UPDATE task_submissions 
			SET status = 'rejected',
			    rejection_reason = $1,
			    reviewed_at = NOW()
			WHERE task_id = $2 AND submission_type = $3 AND status = 'pending'
		`, req.RejectionReason, taskID, req.Milestone)
		if err != nil {
			return response.Error(500, fmt.Sprintf("Failed to update submission: %v", err))
		}

		// Update task status back to previous state
		_, err = pool.Exec(ctx, `
			UPDATE tasks 
			SET status = $1, updated_at = NOW()
			WHERE task_id = $2
		`, rejectedStatus, taskID)
		if err != nil {
			return response.Error(500, fmt.Sprintf("Failed to update task: %v", err))
		}

		return response.Success(map[string]interface{}{
			"message": "Work rejected",
			"status":  rejectedStatus,
		})
	}

	// Handle approval - continue with payment
	newStatus := approvedStatus

	// Calculate payment amount
	rewardFloat := new(big.Float)
	rewardFloat.SetString(task.RewardAmount)
	multiplier := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	rewardFloat.Mul(rewardFloat, multiplier)
	rewardWei, _ := rewardFloat.Int(nil)

	paymentWei := new(big.Int).Mul(rewardWei, big.NewInt(int64(paymentPercent)))
	paymentWei.Div(paymentWei, big.NewInt(10000))

	// Pay milestone on blockchain
	txHash, err := client.PayMilestone(uint64(task.ContractTaskID), paymentWei)
	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to pay milestone: %v", err))
	}

	// Calculate new paid amount
	paidFloat := new(big.Float)
	paidFloat.SetString(task.PaidAmount)
	paymentFloat := new(big.Float).SetInt(paymentWei)
	paymentFloat.Quo(paymentFloat, multiplier)
	paidFloat.Add(paidFloat, paymentFloat)
	newPaidAmount := paidFloat.Text('f', 8)

	// Update submission status to approved
	_, err = pool.Exec(ctx, `
		UPDATE task_submissions 
		SET status = 'approved',
		    reviewed_at = NOW()
		WHERE task_id = $1 AND submission_type = $2 AND status = 'pending'
	`, taskID, req.Milestone)
	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to update submission: %v", err))
	}

	// Update task status and paid amount
	_, err = pool.Exec(ctx, `
		UPDATE tasks 
		SET status = $1, paid_amount = $2, updated_at = CURRENT_TIMESTAMP
		WHERE task_id = $3
	`, newStatus, newPaidAmount, taskID)
	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to update task: %v", err))
	}

	// If completed, update user stats
	if newStatus == models.TaskStatusCompleted {
		_, err = pool.Exec(ctx, `
			UPDATE users 
			SET tasks_completed = tasks_completed + 1,
			    credit_score = credit_score + 100
			WHERE did = (SELECT executor_did FROM tasks WHERE task_id = $1)
		`, taskID)
		if err != nil {
			return response.Error(500, fmt.Sprintf("Failed to update user stats: %v", err))
		}
	}

	return response.Success(ApproveWorkResponse{
		Status: newStatus,
		Payment: PaymentDetail{
			Amount: paymentFloat.Text('f', 8),
			TxHash: txHash,
		},
	})
}

func main() {
	lambda.Start(handler)
}
