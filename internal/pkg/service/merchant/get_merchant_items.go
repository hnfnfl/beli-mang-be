package merchant

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *MerchantService) GetMerchantItems(ctx *gin.Context, data dto.GetMerchantItemsRequest) (*dto.GetMerchantItemsResponse, errs.Response) {
	db := s.db
	var (
		stmt  strings.Builder
		items dto.GetMerchantItemsResponse
		count *int
	)

	// set default value
	items.Data = make([]model.MerchantItem, 0)

	stmt.WriteString("WITH filtered AS (SELECT item_id, name, product_categories, price, image_url, created_at FROM merchant_items WHERE 1=1 ")

	stmt.WriteString(fmt.Sprintf("AND merchant_id = '%s' ", data.MerchantId))

	if data.ItemId != "" {
		stmt.WriteString(fmt.Sprintf("AND item_id = '%s' ", data.ItemId))
	}

	if data.Name != "" {
		stmt.WriteString(fmt.Sprintf("AND name LIKE '%%%s%%' ", data.Name))
	}

	if data.ProductCategory == "<invalid>" {
		return &items, errs.Response{}
	} else if data.ProductCategory != "" {
		stmt.WriteString(fmt.Sprintf("AND product_categories = '%s' ", data.ProductCategory))
	}

	stmt.WriteString(") SELECT(SELECT COUNT(*) FROM filtered) AS total, f.* FROM filtered f ")

	if data.CreatedAt == "ASC" {
		stmt.WriteString("ORDER BY created_at ASC")
	} else {
		stmt.WriteString("ORDER BY created_at DESC")
	}

	stmt.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, data.Offset))

	rows, err := db.Query(ctx, stmt.String())
	if err != nil {
		return nil, errs.NewInternalError("Failed to get merchant items", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			item      model.MerchantItem
			createdAt time.Time
		)
		if err := rows.Scan(
			&count,
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

		items.Data = append(items.Data, item)
	}

	if len(items.Data) == 0 {
		return &items, errs.Response{}
	} else {
		items.Meta = &errs.Meta{
			Limit:  data.Limit,
			Offset: data.Offset,
			Total:  *count,
		}

		return &items, errs.Response{}
	}
}
