package service

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/errs"
)

func (s *Service) InsertMerchant(data model.Merchant) (string, errs.Response) {
	db := s.DB()

	stmt := `INSERT INTO public.merchants (merchant_id, "name", merchant_categories, long, lat, image_url) VALUES($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(stmt, data.MerchantId, data.Name, data.MerchantCategory, data.Long, data.Lat, data.ImageUrl)
	if err != nil {
		return "", errs.NewInternalError("Failed to insert merchant", err)
	}

	return data.MerchantId, errs.Response{}
}
