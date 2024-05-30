package merchant

import (
	"beli-mang/internal/pkg/dto"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MerchantService struct {
	db *pgxpool.Pool
}

type MerchantServiceInterface interface {
	InsertMerchant(ctx *gin.Context, data dto.AddMerchantRequest) *dto.AddMerchantResponse
	GetMerchants(ctx *gin.Context, data dto.GetMerchantsRequest) *dto.GetMerchantsResponse
	InsertMerchantItem(ctx *gin.Context, data dto.AddMerchantItemRequest) *dto.AddMerchantItemResponse
	GetMerchantItems(ctx *gin.Context, data dto.GetMerchantItemsRequest) *dto.GetMerchantItemsResponse
}

func NewMerchantService(db *pgxpool.Pool) *MerchantService {
	return &MerchantService{db}
}

var (
	_ MerchantServiceInterface = &MerchantService{}
)
