package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Host                string
	Port                int
	User                string
	Password            string
	DBName              string
	SSLMode             string        // disable, allow, prefer, require, verify-ca, verify-full
	PoolMaxConns        int           // def: 10
	PoolMinConns        int           // def: 2
	PoolMaxConnIdle     time.Duration // def: 30m
	PoolMaxConnLifetime time.Duration // def: 1h
	ConnTimeout         time.Duration // def: 5s

}

func NewPostgresDB(ctx context.Context, cfg DBConfig) (*pgxpool.Pool, error) {
	if cfg.PoolMaxConns == 0 {
		cfg.PoolMaxConns = 10
	}
	if cfg.PoolMinConns == 0 {
		cfg.PoolMinConns = 2
	}
	if cfg.PoolMaxConnLifetime == 0 {
		cfg.PoolMaxConnLifetime = time.Hour
	}
	if cfg.ConnTimeout == 0 {
		cfg.ConnTimeout = 5 * time.Second
	}
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.PoolMaxConns)
	poolConfig.MinConns = int32(cfg.PoolMinConns)
	poolConfig.MaxConnIdleTime = cfg.PoolMaxConnIdle
	poolConfig.MaxConnLifetime = cfg.PoolMaxConnLifetime
	poolConfig.ConnConfig.ConnectTimeout = cfg.ConnTimeout

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, "SET TIME ZONE'UTC")
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.ConnTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return pool, nil
}

func HealthCheck(ctx context.Context, pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var result int
	err := pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("unexpected health check result: %d", result)
	}

	return nil
}

func CloseDB(pool *pgxpool.Pool) {
	pool.Close()
}
