package service

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"fmt"
	"strings"
	"time"
)

func (s *Service) InsertMerchant(data dto.AddMerchantRequest) (*dto.AddMerchantResponse, errs.Response) {
	db := s.DB()
	var merchant model.Merchant

	stmt := `INSERT INTO merchants (merchant_id, "name", merchant_categories, long, lat, image_url) VALUES($1, $2, $3, $4, $5, $6) RETURNING merchant_id`

	merchantId := db.QueryRow(stmt, data.MerchantId, data.Name, data.MerchantCategory, data.Location.Long, data.Location.Lat, data.ImageUrl)
	if err := merchantId.Scan(&merchant.MerchantId); err != nil {
		return nil, errs.NewInternalError("Failed to insert merchant", err)
	}

	return &dto.AddMerchantResponse{
			MerchantId: merchant.MerchantId,
		},
		errs.Response{}
}

func (s *Service) GetMerchants(data dto.GetMerchantsRequest) (*dto.GetMerchantsResponse, errs.Response) {
	db := s.DB()
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

	rows, err := db.Query(stmt.String())
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

	if merchants.Data == nil {
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

func (s *Service) InsertMerchantItem(data dto.AddMerchantItemRequest) (*dto.AddMerchantItemResponse, errs.Response) {
	db := s.DB()
	var item model.MerchantItem

	// check if merchant exists
	stmt := `SELECT COUNT(*) FROM merchants WHERE merchant_id = $1`
	var count int
	if err := db.QueryRow(stmt, data.MerchantId).Scan(&count); err != nil {
		return nil, errs.NewInternalError("Failed to check merchant", err)
	}

	if count != 0 {
		stmt = `INSERT INTO merchant_items (item_id, merchant_id, "name", price, product_categories, image_url) VALUES($1, $2, $3, $4, $5, $6) RETURNING item_id`

		itemId := db.QueryRow(stmt, data.ItemId, data.MerchantId, data.Name, data.Price, data.ProductCategory, data.ImageUrl)
		if err := itemId.Scan(&item.ItemId); err != nil {
			return nil, errs.NewInternalError("Failed to insert merchant item", err)
		}

	} else {
		return nil, errs.NewNotFoundError(errs.ErrMerchantNotFound)
	}

	return &dto.AddMerchantItemResponse{
			ItemId: item.ItemId,
		},
		errs.Response{}
}

func (s *Service) GetMerchantItems(data dto.GetMerchantItemsRequest) (*dto.GetMerchantItemsResponse, errs.Response) {
	db := s.DB()
	var (
		stmt  strings.Builder
		items dto.GetMerchantItemsResponse
		count *int
	)

	// set default value
	items.Data = make([]model.MerchantItem, 0)

	stmt.WriteString("WITH filtered AS (SELECT item_id, name, product_categories, price, image_url, created_at FROM merchant_items WHERE 1=1 ")

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

	rows, err := db.Query(stmt.String())
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

	if items.Data == nil {
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
