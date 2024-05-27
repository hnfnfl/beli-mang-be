package handler

import (
	"beli-mang/internal/pkg/service"

	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	service *service.Service
}

func NewPurchaseHandler(s *service.Service) *PurchaseHandler {
	return &PurchaseHandler{s}
}

func (h *PurchaseHandler) NearbyMerchant(ctx *gin.Context) {
	
}
