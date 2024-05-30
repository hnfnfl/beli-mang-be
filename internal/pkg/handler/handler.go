package handler

import (
	"beli-mang/internal/pkg/configuration"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	cfg    *configuration.Configuration
	db     *pgxpool.Pool
	logger *logrus.Logger
}

func NewHandler(cfg *configuration.Configuration, db *pgxpool.Pool, logger *logrus.Logger) *Handler {
	return &Handler{
		cfg:    cfg,
		db:     db,
		logger: logger,
	}
}

func (h *Handler) Config() *configuration.Configuration {
	return h.cfg
}

func (h *Handler) DB() *pgxpool.Pool {
	return h.db
}

func (h *Handler) Logger() *logrus.Logger {
	return h.logger
}
