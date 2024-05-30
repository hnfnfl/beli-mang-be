package merchant

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MerchantService struct {
	db *pgxpool.Pool
}

type MerchantServiceInterface interface {
	InsertMerchant(ctx *gin.Context, data dto.AddMerchantRequest) (*dto.AddMerchantResponse, errs.Response)
	GetMerchants(ctx *gin.Context, data dto.GetMerchantsRequest) (*dto.GetMerchantsResponse, errs.Response)
	InsertMerchantItem(ctx *gin.Context, data dto.AddMerchantItemRequest) (*dto.AddMerchantItemResponse, errs.Response)
	GetMerchantItems(ctx *gin.Context, data dto.GetMerchantItemsRequest) (*dto.GetMerchantItemsResponse, errs.Response)
}

func NewMerchantService(db *pgxpool.Pool) *MerchantService {
	return &MerchantService{db}
}

var (
	_ MerchantServiceInterface = &MerchantService{}
)
