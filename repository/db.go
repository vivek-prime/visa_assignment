package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"visa_assignment/config"
)

func InitDB(database config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", database.ConnectionString)
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
