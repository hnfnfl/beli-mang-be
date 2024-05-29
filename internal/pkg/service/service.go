package service

import (
	"beli-mang/internal/pkg/configuration"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	cfg *configuration.Configuration
	db  *pgxpool.Pool
}

func NewService(cfg *configuration.Configuration, db *pgxpool.Pool) *Service {
	return &Service{
		cfg: cfg,
		db:  db,
	}
}

func (s *Service) DB() *pgxpool.Pool {
	return s.db
}

func (s *Service) Config() *configuration.Configuration {
	return s.cfg
}
