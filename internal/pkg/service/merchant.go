package service

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
)

func (s *Service) InsertMerchant(data dto.AddMerchantRequest) (string, errs.Response) {
	db := s.DB()

	stmt := `INSERT INTO merchants (merchant_id, "name", merchant_categories, long, lat, image_url) VALUES($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(stmt, data.MerchantId, data.Name, data.MerchantCategory, data.Location.Long, data.Location.Lat, data.ImageUrl)
	if err != nil {
		return "", errs.NewInternalError("Failed to insert merchant", err)
	}

	return data.MerchantId, errs.Response{}
}

func (s *Service) InsertMerchantItem(data dto.AddMerchantItemRequest) (string, errs.Response) {
	db := s.DB()

	// check if merchant exists
	stmt := `SELECT COUNT(*) FROM merchants WHERE merchant_id = $1`
	var count int
	if err := db.QueryRow(stmt, data.MerchantId).Scan(&count); err != nil {
		return "", errs.NewInternalError("Failed to check merchant", err)
	}

	if count != 0 {
		stmt = `INSERT INTO merchant_items (item_id, merchant_id, "name", price, product_categories, image_url) VALUES($1, $2, $3, $4, $5, $6)`

		_, err := db.Exec(stmt, data.ItemId, data.MerchantId, data.Name, data.Price, data.ProductCategory, data.ImageUrl)
		if err != nil {
			return "", errs.NewInternalError("Failed to insert merchant item", err)
		}
	} else {
		return "", errs.NewNotFoundError(errs.ErrMerchantNotFound)
	}

	return data.ItemId, errs.Response{}
}
