package order

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderService struct {
	db *pgxpool.Pool
}

type OrderServiceInterface interface {
	GetNearbyMerchants(ctx *gin.Context, data dto.GetNearbyMerchantsRequest) (*dto.GetNearbyMerchantsResponse, errs.Response)
	EstimateOrder(ctx *gin.Context, data dto.OrderEstimateRequest) (*dto.OrderEstimateResponse, errs.Response)
}

func NewOrderService(db *pgxpool.Pool) *OrderService {
	return &OrderService{db}
}

var (
	_ OrderServiceInterface = &OrderService{}
)
