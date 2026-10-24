package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

func NewSQLServerDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlserver connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping sqlserver: %w", err)
	}

	return db, nil
}

// InitSchema reads and executes the DDL script (e.g., resources/schema.sql)
func InitSchema(db *sql.DB, schemaFilePath string) error {
	script, err := os.ReadFile(schemaFilePath)
	if err != nil {
		return fmt.Errorf("failed to read schema file %s: %w", schemaFilePath, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if _, err := db.ExecContext(ctx, string(script)); err != nil {
		return fmt.Errorf("failed to execute schema script: %w", err)
	}

	return nil
}
