package handler

import (
	"beli-mang/internal/pkg/service"
)

type PurchaseHandler struct {
	service *service.Service
}

func NewPurchaseHandler(s *service.Service) *PurchaseHandler {
	return &PurchaseHandler{s}
}
