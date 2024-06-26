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

func (s *OrderService) GetNearbyMerchants(ctx *gin.Context, data dto.GetNearbyMerchantsRequest) *dto.GetNearbyMerchantsResponse {
	db := s.db
	var (
		stmt      strings.Builder
		result    dto.GetNearbyMerchantsResponse
		totalData int
	)

	// set default value
	result.Data = make([]dto.GetNearbyMerchantsResponseData, 0)

	userLat := data.Lat
	userLong := data.Long

	// get nearby merchants
	stmt.WriteString("WITH totalCount AS (SELECT COUNT(*) as total FROM merchants)")
	stmt.WriteString(fmt.Sprintf(
		`SELECT *,
			(acos(
				cos(radians(%f)) * cos(radians(lat)) *
				cos(radians(long) - radians(%f)) +
				sin(radians(%f)) * sin(radians(lat))
			)) as distance
		FROM merchants m, totalCount tc WHERE 1=1`,
		userLat, userLong, userLat,
	))

	if data.MerchantId != "" {
		stmt.WriteString(fmt.Sprintf(" AND merchant_id = '%s' ", data.MerchantId))
	}

	if data.Name != "" {
		stmt.WriteString(fmt.Sprintf(" AND name LIKE '%%%s%%' ", data.Name))
	}

	if data.MerchantCategory == "<invalid>" {
		return &result
	} else if data.MerchantCategory != "" {
		stmt.WriteString(fmt.Sprintf(" AND merchant_categories = '%s' ", data.MerchantCategory))
	}

	stmt.WriteString(" ORDER BY distance")
	stmt.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, data.Offset))
	fmt.Println(stmt.String())

	rows, err := db.Query(ctx, stmt.String())
	if err != nil {
		errs.NewInternalError(ctx, "Failed to get nearby merchants", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var (
			merchant  model.Merchant
			items     = make([]model.MerchantItem, 0)
			createdAt time.Time
			distance  float64
			total     int
		)

		if err := rows.Scan(
			&merchant.MerchantId,
			&merchant.Name,
			&merchant.MerchantCategory,
			&merchant.Location.Long,
			&merchant.Location.Lat,
			&merchant.ImageUrl,
			&createdAt,
			&total,
			&distance,
		); err != nil {
			errs.NewInternalError(ctx, "Failed to scan nearby merchants", err)
			return nil
		}

		merchant.CreatedAt = createdAt.Format(time.RFC3339Nano)
		totalData = total

		// get merchant items if data.Name != ""
		if data.Name != "" {
			stmt.Reset()
			stmt.WriteString("WITH totalCount AS (SELECT COUNT(*) as total FROM merchant_items)")
			stmt.WriteString(fmt.Sprintf("SELECT item_id, name, product_categories, price, image_url, created_at, tc.total FROM merchant_items mi, totalCount tc WHERE 1=1 AND merchant_id = '%s' ", merchant.MerchantId))

			stmt.WriteString(fmt.Sprintf("AND name LIKE '%%%s%%' ", data.Name))

			stmt.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, data.Offset))

			itemRows, err := db.Query(ctx, stmt.String())
			if err != nil {
				errs.NewInternalError(ctx, "Failed to get merchant items", err)
				return nil
			}

			var totalItems int
			for itemRows.Next() {
				var (
					item      model.MerchantItem
					createdAt time.Time
				)

				if err := itemRows.Scan(
					&item.ItemId,
					&item.Name,
					&item.ProductCategory,
					&item.Price,
					&item.ImageUrl,
					&createdAt,
					&total,
				); err != nil {
					errs.NewInternalError(ctx, "Failed to scan merchant items", err)
					return nil
				}

				item.CreatedAt = createdAt.Format(time.RFC3339Nano)

				items = append(items, item)
			}
			totalData += totalItems
		}

		result.Data = append(result.Data, dto.GetNearbyMerchantsResponseData{
			Merchant: merchant,
			Items:    items,
		})
	}

	if len(result.Data) == 0 {
		return &result
	}

	result.Meta = &errs.Meta{
		Limit:  data.Limit,
		Offset: data.Offset,
		Total:  totalData,
	}

	return &result
}
