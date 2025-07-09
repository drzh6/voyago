package handler

import (
	"api/voyago/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool             *pgxpool.Pool
	cfg              config.Config
	AccessTokenName  string
	RefreshTokenName string "refresh_token"
}

func NewService(pool *pgxpool.Pool, cfg config.Config) *Service {
	return &Service{
		pool:             pool,
		cfg:              cfg,
		AccessTokenName:  "access_token",
		RefreshTokenName: "refresh_token",
	}
}
