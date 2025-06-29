package handler

import (
	"api/voyago/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool *pgxpool.Pool
	cfg  config.Config
}

func NewService(pool *pgxpool.Pool, cfg config.Config) *Service {
	return &Service{
		pool: pool,
		cfg:  cfg,
	}
}

/*func GetResultFromRows(rows *pgx.Rows) (string, error) {
	var result string
	err := rows.Scan(&result)
}*/
