package order

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *OrderService) PostOrder(ctx *gin.Context, data dto.PostOrderRequest) *dto.PostOrderResponse {

	cache.RLock()
	entry, found := cache.Data[data.CalculatedEstimateId]
	cache.RUnlock()

	if !found {
		errs.NewNotFoundError(ctx, errs.ErrCalculatedEstimateId)
		return nil
	}

	orderID := util.UuidGenerator("ord", 15)
	db := s.db

	userID := ctx.Value("username").(string)

	for _, order := range entry.Request.Orders {
		for _, item := range order.Items {
			query := `
				INSERT INTO order_product (order_id, user_id, merchant_id, item_id, created_at)
				VALUES ($1, $2, $3, $4, $5)
			`
			_, err := db.Exec(ctx, query, orderID, userID, order.MerchantId, item.ItemId, time.Now())
			if err != nil {
				errs.NewInternalError(ctx, "Failed to insert orders", err)
				return nil
			}
		}
	}

	response := dto.PostOrderResponse{
		OrderId: orderID,
	}

	return &response

}
