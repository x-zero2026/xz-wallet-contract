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

type BalanceResponse struct {
	DID        string `json:"did"`
	EthAddress string `json:"eth_address"`
	XZTBalance string `json:"xzt_balance"`
	Username   string `json:"username"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Handle OPTIONS for CORS
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type,Authorization",
				"Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
			},
			Body: "",
		}, nil
	}

	// Validate JWT token
	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		authHeader = request.Headers["authorization"]
	}
	if authHeader == "" {
		return response.Error(401, "Missing authorization header")
	}

	claims, err := auth.ValidateToken(authHeader)
	if err != nil {
		return response.Error(401, fmt.Sprintf("Invalid token: %v", err))
	}

	// Initialize database
	if err := db.InitDB(); err != nil {
		return response.Error(500, fmt.Sprintf("Database error: %v", err))
	}

	// Initialize blockchain client
	client, err := blockchain.InitClient()
	if err != nil {
		return response.Error(500, fmt.Sprintf("Blockchain error: %v", err))
	}

	// Get user info from database
	pool := db.GetPool()
	var ethAddress, username string
	err = pool.QueryRow(ctx, 
		"SELECT eth_address, username FROM users WHERE did = $1",
		claims.DID,
	).Scan(&ethAddress, &username)
	if err != nil {
		return response.Error(404, "User not found")
	}

	// Get balance from blockchain
	balance, err := client.GetBalance(ethAddress)
	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to get balance: %v", err))
	}

	// Convert to decimal (18 decimals)
	balanceFloat := new(big.Float).SetInt(balance)
	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	balanceFloat.Quo(balanceFloat, divisor)
	balanceStr := balanceFloat.Text('f', 8)

	return response.Success(BalanceResponse{
		DID:        claims.DID,
		EthAddress: ethAddress,
		XZTBalance: balanceStr,
		Username:   username,
	})
}

func main() {
	lambda.Start(handler)
}
