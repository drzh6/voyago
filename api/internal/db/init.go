package db

import (
	"api/voyago/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName))
}
