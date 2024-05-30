package order

import (
	"beli-mang/internal/pkg/dto"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderService struct {
	db *pgxpool.Pool
}

type OrderServiceInterface interface {
	GetNearbyMerchants(ctx *gin.Context, data dto.GetNearbyMerchantsRequest) *dto.GetNearbyMerchantsResponse
	EstimateOrder(ctx *gin.Context, data dto.OrderEstimateRequest) *dto.OrderEstimateResponse
}

func NewOrderService(db *pgxpool.Pool) *OrderService {
	return &OrderService{db}
}

var (
	_ OrderServiceInterface = &OrderService{}
)
