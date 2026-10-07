package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLConnectionManager struct {
	db *sql.DB
}

func NewMySQLConnectionManager() *MySQLConnectionManager {
	return &MySQLConnectionManager{}
}

func (m *MySQLConnectionManager) Connect(ctx context.Context, config ConnectionConfig) error {
	var dsn string
	if config.DSN != "" {
		dsn = config.DSN
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci&readTimeout=30s&writeTimeout=30s",
			config.User, config.Password, config.Host, config.Port, config.Database)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		err = db.Close()
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to ping database: %w", err)
	}

	m.db = db
	return nil
}

func (m *MySQLConnectionManager) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

func (m *MySQLConnectionManager) Ping(ctx context.Context) error {
	if m.db == nil {
		return fmt.Errorf("not connected")
	}
	return m.db.PingContext(ctx)
}

func (m *MySQLConnectionManager) DB() *sql.DB {
	return m.db
}
