package config

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDatabase(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
    poolConfig, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, err
    }

    poolConfig.MaxConnIdleTime = 5 * time.Minute
    poolConfig.MaxConns = 10

    return pgxpool.NewWithConfig(ctx, poolConfig)
}
