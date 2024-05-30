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

func (s *MerchantService) GetMerchants(ctx *gin.Context, data dto.GetMerchantsRequest) (*dto.GetMerchantsResponse, errs.Response) {
	db := s.db
	var (
		stmt      strings.Builder
		merchants dto.GetMerchantsResponse
		count     *int
	)

	// set default value
	merchants.Data = make([]model.Merchant, 0)

	stmt.WriteString("WITH filtered AS (SELECT * FROM merchants WHERE 1=1 ")

	if data.MerchantId != "" {
		stmt.WriteString(fmt.Sprintf("AND merchant_id = '%s' ", data.MerchantId))
	}

	if data.Name != "" {
		stmt.WriteString(fmt.Sprintf("AND name LIKE '%%%s%%' ", data.Name))
	}

	if data.MerchantCategory == "<invalid>" {
		return &merchants, errs.Response{}
	} else if data.MerchantCategory != "" {
		stmt.WriteString(fmt.Sprintf("AND merchant_categories = '%s' ", data.MerchantCategory))
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
		return nil, errs.NewInternalError("Failed to get merchants", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			merchant  model.Merchant
			createdAt time.Time
		)
		if err := rows.Scan(
			&count,
			&merchant.MerchantId,
			&merchant.Name,
			&merchant.MerchantCategory,
			&merchant.Location.Long,
			&merchant.Location.Lat,
			&merchant.ImageUrl,
			&createdAt,
		); err != nil {
			return nil, errs.NewInternalError("Failed to scan merchants", err)
		}

		merchant.CreatedAt = createdAt.Format(time.RFC3339Nano)

		merchants.Data = append(merchants.Data, merchant)
	}

	if len(merchants.Data) == 0 {
		return &merchants, errs.Response{}
	} else {
		merchants.Meta = &errs.Meta{
			Limit:  data.Limit,
			Offset: data.Offset,
			Total:  *count,
		}

		return &merchants, errs.Response{}
	}
}
