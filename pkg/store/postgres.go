package store

import (
	"context"
	"fmt"
	"time"

	"github.com/d4vz/rinha-de-backend-2025/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DB *pgxpool.Pool
)

func ConnectDB() error {
	dbHost := config.GetEnvOrDefault("DB_HOST", "localhost")
	dbPort := config.GetEnvOrDefaultInt("DB_PORT", 5432)
	dbUser := config.GetEnvOrDefault("DB_USER", "rinha")
	dbPassword := config.GetEnvOrDefault("DB_PASSWORD", "rinha")
	dbName := config.GetEnvOrDefault("DB_NAME", "rinha")
	dbMaxConns := config.GetEnvOrDefaultInt("DB_MAX_CONNS", 10)
	dbMinConns := config.GetEnvOrDefaultInt("DB_MIN_CONNS", 2)

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	dbpool, err := pgxpool.New(context.Background(), connString)

	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	DB = dbpool
	DB.Config().MaxConnIdleTime = 10 * time.Minute
	DB.Config().MaxConnLifetime = 2 * time.Hour
	DB.Config().MaxConns = int32(dbMaxConns)
	DB.Config().MinConns = int32(dbMinConns)
	DB.Config().HealthCheckPeriod = 10 * time.Minute

	return nil
}

func GetDB() *pgxpool.Pool {
	return DB
}

func MigrateDB() error {
	_, err := DB.Exec(context.Background(), `
		CREATE UNLOGGED TABLE IF NOT EXISTS payments (
			id SERIAL PRIMARY KEY,
			correlation_id UUID NOT NULL,
			amount NUMERIC(10, 2) NOT NULL,
			processor VARCHAR(10) NOT NULL,
			processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)

	if err != nil {
		return fmt.Errorf("failed to create payments table: %w", err)
	}

	return nil
}
