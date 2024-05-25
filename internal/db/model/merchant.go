package model

import "time"

type Merchant struct {
	MerchantId       string    `db:"merchant_id"`
	Name             string    `db:"name"`
	MerchantCategory string    `db:"merchant_categories"`
	Lat              float64   `db:"lat"`
	Long             float64   `db:"long"`
	ImageUrl         string    `db:"image_url"`
	CreatedAt        time.Time `db:"created_at"`
}

type MerchantItems struct {
	ItemId          string    `db:"item_id"`
	MerchantId      string    `db:"merchant_id"`
	Name            string    `db:"name"`
	ProductCategory string    `db:"product_categories"`
	Price           float64   `db:"price"`
	ImageUrl        string    `db:"image_url"`
	CreatedAt       time.Time `db:"created_at"`
}
