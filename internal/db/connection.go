package db

import (
	"context"
	"database/sql"
)

// ConnectionManager defines the interface for database connection lifecycle management.
// This separates connection concerns from repository operations.
type ConnectionManager interface {
	// Connect establishes a database connection with the given configuration.
	Connect(ctx context.Context, config ConnectionConfig) error

	// Close closes the database connection.
	Close() error

	// Ping verifies the database connection is alive.
	Ping(ctx context.Context) error

	// DB returns the underlying sql.DB for direct query execution.
	DB() *sql.DB
}

// ConnectionConfig holds the configuration for database connection.
type ConnectionConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	DSN      string
}
