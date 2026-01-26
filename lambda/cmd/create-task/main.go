package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/x-zero/xz-wallet/pkg/auth"
	"github.com/x-zero/xz-wallet/pkg/blockchain"
	"github.com/x-zero/xz-wallet/pkg/db"
	"github.com/x-zero/xz-wallet/pkg/response"
)

type CreateTaskRequest struct {
	ProjectID          string `json:"project_id"`
	TaskName           string `json:"task_name"`
	TaskDescription    string `json:"task_description"`
	AcceptanceCriteria string `json:"acceptance_criteria"`
	RewardAmount       string `json:"reward_amount"`
	Visibility         string `json:"visibility"`
}

type CreateTaskResponse struct {
	TaskID         string `json:"task_id"`
	ContractTaskID int64  `json:"contract_task_id"`
	TxHash         string `json:"tx_hash"`
	Status         string `json:"status"`
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

	// Parse request
	var req CreateTaskRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.Error(400, "Invalid request body")
	}

	// Validate input
	if req.ProjectID == "" || req.TaskName == "" || req.RewardAmount == "" {
		return response.Error(400, "Missing required fields")
	}
	if req.Visibility != "project" && req.Visibility != "global" {
		return response.Error(400, "Invalid visibility")
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

	// Get user's eth_address
	var ethAddress string
	err = pool.QueryRow(ctx, "SELECT eth_address FROM users WHERE did = $1", claims.DID).Scan(&ethAddress)
	if err != nil {
		return response.Error(404, "User not found")
	}

	// Convert amount to wei
	amountFloat := new(big.Float)
	amountFloat.SetString(req.RewardAmount)
	multiplier := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	amountFloat.Mul(amountFloat, multiplier)
	amountWei, _ := amountFloat.Int(nil)

	// Check if user has approved escrow contract
	allowance, err := client.Token.Allowance(&bind.CallOpts{}, common.HexToAddress(ethAddress), client.EscrowAddress)
	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to check allowance: %v", err))
	}

	// If allowance is insufficient, call approve-escrow Lambda
	if allowance.Cmp(amountWei) < 0 {
		// Check user balance first
		userBalance, err := client.Token.BalanceOf(&bind.CallOpts{}, common.HexToAddress(ethAddress))
		if err != nil {
			return response.Error(500, fmt.Sprintf("Failed to check user balance: %v", err))
		}
		
		if userBalance.Cmp(amountWei) < 0 {
			return response.Error(400, fmt.Sprintf("Insufficient XZT balance. Required: %s XZT, Available: %s XZT", 
				req.RewardAmount, 
				new(big.Float).Quo(new(big.Float).SetInt(userBalance), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))).Text('f', 2)))
		}
		
		// Call approve-escrow Lambda to automatically approve
		fmt.Printf("Allowance insufficient, calling approve-escrow Lambda...\n")
		err = callApproveEscrow(authHeader, client.TokenAddress.Hex(), client.EscrowAddress.Hex())
		if err != nil {
			return response.Error(500, fmt.Sprintf("Failed to approve escrow contract: %v. Please try again.", err))
		}
		fmt.Printf("Escrow approved successfully, waiting for network propagation...\n")
		
		// Wait a bit for the approval transaction to fully propagate
		// This helps avoid "in-flight transaction limit" errors on Alchemy free tier
		time.Sleep(3 * time.Second)
	}

	// Insert into database FIRST with pending status
	// This way if blockchain fails, we can mark as cancelled
	var taskID string
	err = pool.QueryRow(ctx, `
		INSERT INTO tasks (
			contract_task_id, project_id, creator_did, task_name, 
			task_description, acceptance_criteria, reward_amount, 
			visibility, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING task_id
	`, -1, req.ProjectID, claims.DID, req.TaskName,
		req.TaskDescription, req.AcceptanceCriteria, req.RewardAmount,
		req.Visibility, "pending").Scan(&taskID)
	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to save task: %v", err))
	}
	fmt.Printf("Task saved to database with ID: %s, now creating on blockchain...\n", taskID)

	// Create task on blockchain
	contractTaskID, txHash, err := client.CreateTask(ethAddress, amountWei)
	if err != nil {
		// Blockchain failed, mark task as cancelled in database
		_, updateErr := pool.Exec(ctx, `
			UPDATE tasks 
			SET status = 'cancelled', 
			    updated_at = NOW(),
			    task_description = task_description || E'\n\n[系统消息] 区块链创建失败: ' || $1
			WHERE task_id = $2
		`, err.Error(), taskID)
		if updateErr != nil {
			fmt.Printf("Failed to update task status: %v\n", updateErr)
		}
		return response.Error(500, fmt.Sprintf("Failed to create task on blockchain: %v", err))
	}

	// Update task with contract_task_id and set status to bidding
	_, err = pool.Exec(ctx, `
		UPDATE tasks 
		SET contract_task_id = $1,
		    status = 'bidding',
		    updated_at = NOW()
		WHERE task_id = $2
	`, contractTaskID, taskID)
	if err != nil {
		// This is bad - blockchain succeeded but database update failed
		// Log the orphaned task for manual recovery
		fmt.Printf("CRITICAL: Task created on blockchain (contract_task_id=%d, tx=%s) but failed to update database (task_id=%s): %v\n", 
			contractTaskID, txHash, taskID, err)
		return response.Error(500, fmt.Sprintf("Task created on blockchain but database update failed. Contract Task ID: %d, TX: %s. Please contact support.", contractTaskID, txHash))
	}

	return response.Success(CreateTaskResponse{
		TaskID:         taskID,
		ContractTaskID: int64(contractTaskID),
		TxHash:         txHash,
		Status:         "bidding",
	})
}

// callApproveEscrow calls the approve-escrow Lambda function in did-login-lambda
func callApproveEscrow(authToken, tokenAddress, spenderAddress string) error {
	// Get DID Login API URL from environment
	didLoginAPIURL := os.Getenv("DID_LOGIN_API_URL")
	if didLoginAPIURL == "" {
		return fmt.Errorf("DID_LOGIN_API_URL not configured")
	}

	// Prepare request body
	reqBody := map[string]string{
		"token_address":   tokenAddress,
		"spender_address": spenderAddress,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := didLoginAPIURL + "/api/approve-escrow"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != 200 {
		return fmt.Errorf("approve-escrow API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if success, ok := result["success"].(bool); !ok || !success {
		errMsg := "unknown error"
		if msg, ok := result["error"].(string); ok {
			errMsg = msg
		}
		return fmt.Errorf("approval failed: %s", errMsg)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
