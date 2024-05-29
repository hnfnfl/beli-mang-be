package handler

import (
	"beli-mang/internal/pkg/configuration"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	cfg    *configuration.Configuration
	logger *logrus.Logger
}

func NewHandler(cfg *configuration.Configuration, logger *logrus.Logger) *Handler {
	return &Handler{
		cfg:    cfg,
		logger: logger,
	}
}

func (h *Handler) Config() *configuration.Configuration {
	return h.cfg
}

func (h *Handler) Logger() *logrus.Logger {
	return h.logger
}
