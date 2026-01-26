package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// CreateTask creates a new task and locks XZT in escrow
func (c *BlockchainClient) CreateTask(creatorAddress string, amount *big.Int) (uint64, string, error) {
	creator := common.HexToAddress(creatorAddress)
	executor := common.HexToAddress("0x0000000000000000000000000000000000000000") // No executor yet

	// First, approve escrow to spend creator's tokens (if not already approved)
	// Note: In production, this should be done once per user during registration

	// Create task
	tx, err := c.Escrow.CreateTask(c.AdminAuth, creator, executor, amount)
	if err != nil {
		return 0, "", fmt.Errorf("failed to create task: %w", err)
	}

	// Wait for transaction
	receipt, err := bind.WaitMined(context.Background(), c.Client, tx)
	if err != nil {
		return 0, "", fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status == 0 {
		return 0, "", fmt.Errorf("transaction failed")
	}

	// Parse TaskCreated event to get task ID
	// For simplicity, we'll return the next task ID - 1
	// In production, parse the event logs
	nextTaskID, err := c.Escrow.NextTaskId(&bind.CallOpts{})
	if err != nil {
		return 0, "", fmt.Errorf("failed to get next task ID: %w", err)
	}

	taskID := nextTaskID.Uint64() - 1

	return taskID, tx.Hash().Hex(), nil
}

// SetExecutor sets the executor for a task
func (c *BlockchainClient) SetExecutor(taskID uint64, executorAddress string) (string, error) {
	executor := common.HexToAddress(executorAddress)
	taskIDBig := big.NewInt(int64(taskID))

	tx, err := c.Escrow.SetExecutor(c.AdminAuth, taskIDBig, executor)
	if err != nil {
		return "", fmt.Errorf("failed to set executor: %w", err)
	}

	receipt, err := bind.WaitMined(context.Background(), c.Client, tx)
	if err != nil {
		return "", fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status == 0 {
		return "", fmt.Errorf("transaction failed")
	}

	return tx.Hash().Hex(), nil
}

// PayMilestone pays a milestone to the executor
func (c *BlockchainClient) PayMilestone(taskID uint64, amount *big.Int) (string, error) {
	taskIDBig := big.NewInt(int64(taskID))

	tx, err := c.Escrow.PayMilestone(c.AdminAuth, taskIDBig, amount)
	if err != nil {
		return "", fmt.Errorf("failed to pay milestone: %w", err)
	}

	receipt, err := bind.WaitMined(context.Background(), c.Client, tx)
	if err != nil {
		return "", fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status == 0 {
		return "", fmt.Errorf("transaction failed")
	}

	return tx.Hash().Hex(), nil
}

// CancelTask cancels a task with refund distribution
func (c *BlockchainClient) CancelTask(taskID uint64, executorAmount *big.Int) (string, error) {
	taskIDBig := big.NewInt(int64(taskID))

	tx, err := c.Escrow.CancelTask(c.AdminAuth, taskIDBig, executorAmount)
	if err != nil {
		return "", fmt.Errorf("failed to cancel task: %w", err)
	}

	receipt, err := bind.WaitMined(context.Background(), c.Client, tx)
	if err != nil {
		return "", fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status == 0 {
		return "", fmt.Errorf("transaction failed")
	}

	return tx.Hash().Hex(), nil
}

// GetTask retrieves task details from contract
func (c *BlockchainClient) GetTask(taskID uint64) (creator, executor string, totalAmount, paidAmount *big.Int, cancelled bool, err error) {
	taskIDBig := big.NewInt(int64(taskID))

	result, err := c.Escrow.GetTask(&bind.CallOpts{}, taskIDBig)
	if err != nil {
		return "", "", nil, nil, false, fmt.Errorf("failed to get task: %w", err)
	}

	return result.Creator.Hex(), result.Executor.Hex(), result.TotalAmount, result.PaidAmount, result.Cancelled, nil
}

// GetRemainingAmount returns the remaining amount for a task
func (c *BlockchainClient) GetRemainingAmount(taskID uint64) (*big.Int, error) {
	taskIDBig := big.NewInt(int64(taskID))

	remaining, err := c.Escrow.GetRemainingAmount(&bind.CallOpts{}, taskIDBig)
	if err != nil {
		return nil, fmt.Errorf("failed to get remaining amount: %w", err)
	}

	return remaining, nil
}
