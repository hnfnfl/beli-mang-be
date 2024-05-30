package merchant

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"

	"github.com/gin-gonic/gin"
)

func (s *MerchantService) InsertMerchant(ctx *gin.Context, data dto.AddMerchantRequest) *dto.AddMerchantResponse {
	db := s.db
	var merchant model.Merchant

	stmt := `INSERT INTO merchants (merchant_id, "name", merchant_categories, long, lat, image_url) VALUES($1, $2, $3, $4, $5, $6) RETURNING merchant_id`

	merchantId := db.QueryRow(ctx, stmt, data.MerchantId, data.Name, data.MerchantCategory, data.Location.Long, data.Location.Lat, data.ImageUrl)
	if err := merchantId.Scan(&merchant.MerchantId); err != nil {
		errs.NewInternalError(ctx, "Failed to insert merchant", err)
		return nil
	}

	return &dto.AddMerchantResponse{
		MerchantId: merchant.MerchantId,
	}
}
