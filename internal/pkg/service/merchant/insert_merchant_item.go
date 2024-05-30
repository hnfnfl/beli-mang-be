package merchant

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"

	"github.com/gin-gonic/gin"
)

func (s *MerchantService) InsertMerchantItem(ctx *gin.Context, data dto.AddMerchantItemRequest) *dto.AddMerchantItemResponse {
	db := s.db
	var item model.MerchantItem

	// check if merchant exists
	stmt := `SELECT COUNT(*) FROM merchants WHERE merchant_id = $1`
	var count int
	if err := db.QueryRow(ctx, stmt, data.MerchantId).Scan(&count); err != nil {
		errs.NewInternalError(ctx, "Failed to check merchant", err)
		return nil
	}

	if count != 0 {
		stmt = `INSERT INTO merchant_items (item_id, merchant_id, "name", price, product_categories, image_url) VALUES($1, $2, $3, $4, $5, $6) RETURNING item_id`

		itemId := db.QueryRow(ctx, stmt, data.ItemId, data.MerchantId, data.Name, data.Price, data.ProductCategory, data.ImageUrl)
		if err := itemId.Scan(&item.ItemId); err != nil {
			errs.NewInternalError(ctx, "Failed to insert merchant item", err)
			return nil
		}

	} else {
		errs.NewNotFoundError(ctx, errs.ErrMerchantNotFound)
		return nil
	}

	return &dto.AddMerchantItemResponse{
		ItemId: item.ItemId,
	}
}
