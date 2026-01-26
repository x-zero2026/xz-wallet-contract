package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

// InitDB initializes the database connection pool (singleton)
func InitDB() error {
	var err error
	once.Do(func() {
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			err = fmt.Errorf("DATABASE_URL environment variable not set")
			return
		}

		// Parse config and disable prepared statements for Lambda
		config, parseErr := pgxpool.ParseConfig(dbURL)
		if parseErr != nil {
			err = fmt.Errorf("failed to parse database URL: %w", parseErr)
			return
		}

		// Disable prepared statements to avoid conflicts in Lambda
		config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			err = fmt.Errorf("failed to create connection pool: %w", err)
			return
		}

		// Test connection
		if err = pool.Ping(context.Background()); err != nil {
			err = fmt.Errorf("failed to ping database: %w", err)
			return
		}
	})
	return err
}

// GetPool returns the database connection pool
func GetPool() *pgxpool.Pool {
	return pool
}

// Close closes the database connection pool
func Close() {
	if pool != nil {
		pool.Close()
	}
}
