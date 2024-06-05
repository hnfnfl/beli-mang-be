package order

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *OrderService) GetOrders(ctx *gin.Context, data dto.GetOrdersRequest) *[]dto.GetOrdersResponse {

	db := s.db
	var (
		stmt strings.Builder
	)
	// set default value
	getOrderResponse := make([]dto.GetOrdersResponse, 0)
	orders := map[string]*dto.GetOrdersResponse{}

	query := `
	SELECT 
		op.order_id, 
		m.merchant_id, 
		m.name, 
		m.merchant_categories, 
		m.image_url, 
		m.lat, 
		m.long, 
		m.created_at,
		i.item_id, 
		i.name, 
		i.product_categories, 
		i.price, 
		op.quantity, 
		i.image_url, 
		i.created_at
	FROM order_product op
	JOIN merchants m ON op.merchant_id = m.merchant_id
	JOIN merchant_items i ON op.item_id = i.item_id
	WHERE 1=1 
	`

	stmt.WriteString(query)

	if data.MerchantId != "" {
		stmt.WriteString(fmt.Sprintf("AND op.merchant_id = '%s' ", data.MerchantId))
	}

	if data.Name != "" {
		stmt.WriteString(fmt.Sprintf("AND (LOWER(m.name) LIKE '%%%s%%' OR LOWER(i.name) LIKE '%%%s%%'", data.Name, data.Name))
	}

	if data.MerchantCategory == "<invalid>" {
		return &getOrderResponse
	} else if data.MerchantCategory != "" {
		stmt.WriteString(fmt.Sprintf("AND m.merchant_categories = '%s' ", data.MerchantCategory))
	}

	stmt.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, data.Offset))

	rows, err := db.Query(ctx, stmt.String())
	if err != nil {
		errs.NewInternalError(ctx, "Failed to get orders", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var (
			merchant                         model.Merchant
			orderId                          string
			item                             dto.OrderItemResponse
			merchantCreatedAt, itemCreatedAt time.Time
		)

		if err := rows.Scan(
			&orderId,
			&merchant.MerchantId,
			&merchant.Name,
			&merchant.MerchantCategory,
			&merchant.ImageUrl,
			&merchant.Location.Lat,
			&merchant.Location.Long,
			&merchantCreatedAt,
			&item.ItemId,
			&item.Name,
			&item.ProductCategory,
			&item.Price,
			&item.Quantity,
			&item.ImageURL,
			&itemCreatedAt,
		); err != nil {
			errs.NewInternalError(ctx, "Failed to scan orders", err)
			return nil
		}

		merchant.CreatedAt = merchantCreatedAt.Format(time.RFC3339Nano)

		if orders[orderId] == nil {
			orders[orderId] = &dto.GetOrdersResponse{
				OrderId: orderId,
				Orders: []dto.DetailOrderResponse{{
					Merchant: merchant,
					Items:    []dto.OrderItemResponse{},
				}},
			}
		}

		item.CreatedAt = itemCreatedAt.Format(time.RFC3339Nano)

		orders[orderId].Orders[0].Items = append(orders[orderId].Orders[0].Items, item)
	}
	for _, order := range orders {
		getOrderResponse = append(getOrderResponse, *order)
	}

	return &getOrderResponse
}
