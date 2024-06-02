package order

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *OrderService) PostOrder(ctx *gin.Context, data dto.PostOrderRequest) *dto.PostOrderResponse {
	cachedData, found := cache.Get(data.CalculatedEstimateId)
	if !found {
		errs.NewNotFoundError(ctx, errs.ErrCalculatedEstimateId)
		return nil
	}

	orderID := util.UuidGenerator("ord", 15)
	db := s.db
	entry := cachedData.(*dto.CacheItem)
	userID := ctx.Value("username").(string)

	for _, order := range entry.Request.Orders {
		for _, item := range order.Items {
			query := `
				INSERT INTO order_product (order_id, user_id, merchant_id, item_id, quantity, created_at)
				VALUES ($1, $2, $3, $4, $5, $6)
			`
			_, err := db.Exec(ctx, query, orderID, userID, order.MerchantId, item.ItemId, item.Quantity, time.Now())
			if err != nil {
				errs.NewInternalError(ctx, "Failed to insert orders", err)
				return nil
			}
		}
	}

	return &dto.PostOrderResponse{
		OrderId: orderID,
	}
}
