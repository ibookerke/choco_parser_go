package pgx

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	connMaxLifetime     = time.Minute * 10
	connMaxIdleLifetime = time.Second * 30
	maxOpenConnections  = 100
	maxIdleConnections  = maxOpenConnections / 3
)

func NewPgxPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool

	err := retry.Do(
		func() error {
			poolConfig, err := pgxpool.ParseConfig(dsn)

			poolConfig.MaxConnIdleTime = connMaxIdleLifetime
			poolConfig.MaxConnLifetime = connMaxLifetime
			poolConfig.MaxConns = maxOpenConnections
			poolConfig.MinConns = maxIdleConnections

			pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
			if err != nil {
				return fmt.Errorf("failed to create database connection pool: %w", err)
			}

			return nil
		},
		retry.Attempts(6),
		retry.Delay(time.Second),
		retry.DelayType(retry.FixedDelay),
		retry.LastErrorOnly(true),
	)

	if err != nil {
		return nil, err
	}

	return pool, nil
}
