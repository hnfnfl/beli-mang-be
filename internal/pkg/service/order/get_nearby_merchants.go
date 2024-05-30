package order

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *OrderService) GetNearbyMerchants(ctx *gin.Context, data dto.GetNearbyMerchantsRequest) (*dto.GetNearbyMerchantsResponse, errs.Response) {
	db := s.db
	var (
		stmt   strings.Builder
		result dto.GetNearbyMerchantsResponse
	)

	// set default value
	result.Data.Merchant = make([]model.Merchant, 0)
	result.Data.Items = make([]model.MerchantItem, 0)

	userLat := data.Lat
	userLong := data.Long

	// get nearby merchants
	stmt.WriteString("SELECT * FROM merchants WHERE 1=1")

	if data.MerchantId != "" {
		stmt.WriteString(fmt.Sprintf("AND merchant_id = '%s' ", data.MerchantId))
	}

	if data.Name != "" {
		stmt.WriteString(fmt.Sprintf("AND name LIKE '%%%s%%' ", data.Name))
	}

	if data.MerchantCategory == "<invalid>" {
		return &result, errs.Response{}
	} else if data.MerchantCategory != "" {
		stmt.WriteString(fmt.Sprintf("AND merchant_categories = '%s' ", data.MerchantCategory))
	}

	stmt.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, data.Offset))

	rows, err := db.Query(ctx, stmt.String())
	if err != nil {
		return nil, errs.NewInternalError("Failed to get nearby merchants", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			merchant  model.Merchant
			createdAt time.Time
		)

		if err := rows.Scan(
			&merchant.MerchantId,
			&merchant.Name,
			&merchant.MerchantCategory,
			&merchant.Location.Long,
			&merchant.Location.Lat,
			&merchant.ImageUrl,
			&createdAt,
		); err != nil {
			return nil, errs.NewInternalError("Failed to scan nearby merchants", err)
		}

		merchant.Distance = util.Haversine(userLat, userLong, merchant.Location.Lat, merchant.Location.Long)
		merchant.CreatedAt = createdAt.Format(time.RFC3339Nano)

		result.Data.Merchant = append(result.Data.Merchant, merchant)
	}

	sort.Slice(result.Data.Merchant, func(i, j int) bool {
		return result.Data.Merchant[i].Distance < result.Data.Merchant[j].Distance
	})

	if len(result.Data.Merchant) == 0 {
		return &result, errs.Response{}
	}

	// get merchant items
	if data.Name != "" {
		stmt.Reset()
		stmt.WriteString("SELECT item_id, name, product_categories, price, image_url, created_at FROM merchant_items WHERE 1=1 ")

		stmt.WriteString(fmt.Sprintf("AND name LIKE '%%%s%%' ", data.Name))

		stmt.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, data.Offset))

		rows, err = db.Query(ctx, stmt.String())
		if err != nil {
			return nil, errs.NewInternalError("Failed to get merchant items", err)
		}

		for rows.Next() {
			var (
				item      model.MerchantItem
				createdAt time.Time
			)
			if err := rows.Scan(
				&item.ItemId,
				&item.Name,
				&item.ProductCategory,
				&item.Price,
				&item.ImageUrl,
				&createdAt,
			); err != nil {
				return nil, errs.NewInternalError("Failed to scan merchant items", err)
			}

			item.CreatedAt = createdAt.Format(time.RFC3339Nano)

			result.Data.Items = append(result.Data.Items, item)
		}
	}

	return &result, errs.Response{}
}
