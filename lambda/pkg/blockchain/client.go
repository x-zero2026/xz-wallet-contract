package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/x-zero/xz-wallet/pkg/blockchain/contracts"
)

var (
	client     *BlockchainClient
	clientOnce sync.Once
)

// BlockchainClient wraps Ethereum client and contract instances
type BlockchainClient struct {
	Client       *ethclient.Client
	Token        *contracts.XZToken
	Escrow       *contracts.TaskEscrow
	AdminAuth    *bind.TransactOpts
	ChainID      *big.Int
	TokenAddress common.Address
	EscrowAddress common.Address
}

// InitClient initializes the blockchain client (singleton)
func InitClient() (*BlockchainClient, error) {
	var err error
	clientOnce.Do(func() {
		// Get environment variables
		rpcURL := os.Getenv("SEPOLIA_RPC_URL")
		if rpcURL == "" {
			err = fmt.Errorf("SEPOLIA_RPC_URL not set")
			return
		}

		tokenAddr := os.Getenv("XZT_TOKEN_ADDRESS")
		if tokenAddr == "" {
			err = fmt.Errorf("XZT_TOKEN_ADDRESS not set")
			return
		}

		escrowAddr := os.Getenv("TASK_ESCROW_ADDRESS")
		if escrowAddr == "" {
			err = fmt.Errorf("TASK_ESCROW_ADDRESS not set")
			return
		}

		adminPrivateKey := os.Getenv("ADMIN_WALLET_PRIVATE_KEY")
		if adminPrivateKey == "" {
			err = fmt.Errorf("ADMIN_WALLET_PRIVATE_KEY not set")
			return
		}

		// Connect to Ethereum client
		ethClient, err := ethclient.Dial(rpcURL)
		if err != nil {
			err = fmt.Errorf("failed to connect to Ethereum client: %w", err)
			return
		}

		// Get chain ID
		chainID, err := ethClient.ChainID(context.Background())
		if err != nil {
			err = fmt.Errorf("failed to get chain ID: %w", err)
			return
		}

		// Parse admin private key
		privateKeyHex := adminPrivateKey
		if len(privateKeyHex) > 2 && privateKeyHex[:2] == "0x" {
			privateKeyHex = privateKeyHex[2:] // Remove "0x" prefix if present
		}
		privateKey, err := crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			err = fmt.Errorf("failed to parse admin private key: %w", err)
			return
		}

		// Create admin auth
		adminAuth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
		if err != nil {
			err = fmt.Errorf("failed to create admin auth: %w", err)
			return
		}

		// Set gas parameters
		adminAuth.GasLimit = 500000 // Default gas limit

		// Parse contract addresses
		tokenAddress := common.HexToAddress(tokenAddr)
		escrowAddress := common.HexToAddress(escrowAddr)

		// Create contract instances
		token, err := contracts.NewXZToken(tokenAddress, ethClient)
		if err != nil {
			err = fmt.Errorf("failed to create token contract instance: %w", err)
			return
		}

		escrow, err := contracts.NewTaskEscrow(escrowAddress, ethClient)
		if err != nil {
			err = fmt.Errorf("failed to create escrow contract instance: %w", err)
			return
		}

		client = &BlockchainClient{
			Client:        ethClient,
			Token:         token,
			Escrow:        escrow,
			AdminAuth:     adminAuth,
			ChainID:       chainID,
			TokenAddress:  tokenAddress,
			EscrowAddress: escrowAddress,
		}
	})

	return client, err
}

// GetClient returns the singleton blockchain client
func GetClient() *BlockchainClient {
	return client
}

// GetBalance returns XZT balance for an address
func (c *BlockchainClient) GetBalance(address string) (*big.Int, error) {
	addr := common.HexToAddress(address)
	balance, err := c.Token.BalanceOf(&bind.CallOpts{}, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return balance, nil
}

// Transfer transfers XZT from admin to recipient
func (c *BlockchainClient) Transfer(to string, amount *big.Int) (string, error) {
	toAddr := common.HexToAddress(to)
	
	tx, err := c.Token.Transfer(c.AdminAuth, toAddr, amount)
	if err != nil {
		return "", fmt.Errorf("failed to transfer: %w", err)
	}

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), c.Client, tx)
	if err != nil {
		return "", fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status == 0 {
		return "", fmt.Errorf("transaction failed")
	}

	return tx.Hash().Hex(), nil
}

// CreateUserAuth creates a TransactOpts for a user (for approval)
func (c *BlockchainClient) CreateUserAuth(privateKey *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, c.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user auth: %w", err)
	}
	auth.GasLimit = 100000
	return auth, nil
}
