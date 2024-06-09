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
		stmt.WriteString(fmt.Sprintf(" AND op.merchant_id = '%s' ", data.MerchantId))
	}

	if data.Name != "" {
		stmt.WriteString(fmt.Sprintf(" AND m.name LIKE '%%%s%%' OR i.name LIKE '%%%s%%'", data.Name, data.Name))
	}

	if data.MerchantCategory == "<invalid>" {
		return &getOrderResponse
	} else if data.MerchantCategory != "" {
		stmt.WriteString(fmt.Sprintf(" AND m.merchant_categories = '%s' ", data.MerchantCategory))
	}

	stmt.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, data.Offset))

	rows, err := db.Query(ctx, stmt.String())
	if err != nil {
		errs.NewInternalError(ctx, "Failed to get orders", err)
		return nil
	}
	defer rows.Close()

	orders := map[string]map[string]*dto.OrderResponse{}

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
		item.CreatedAt = itemCreatedAt.Format(time.RFC3339Nano)

		if orders[orderId] == nil {
			orders[orderId] = make(map[string]*dto.OrderResponse)
		}

		if orders[orderId][merchant.MerchantId] == nil {
			orders[orderId][merchant.MerchantId] = &dto.OrderResponse{
				Merchant: merchant,
				Items:    []dto.OrderItemResponse{},
			}
		}

		orders[orderId][merchant.MerchantId].Items = append(orders[orderId][merchant.MerchantId].Items, item)
	}

	for orderId, order := range orders {
		var or []dto.OrderResponse
		for _, o := range order {
			or = append(or, *o)
		}

		getOrderResponse = append(getOrderResponse, dto.GetOrdersResponse{
			OrderId: orderId,
			Orders:  or,
		})
	}

	return &getOrderResponse
}
