package model

type Merchant struct {
	MerchantId       string   `json:"merchant_id"`
	Name             string   `json:"name"`
	MerchantCategory string   `json:"merchant_categories"`
	ImageUrl         string   `json:"image_url"`
	Location         Location `json:"location"`
	CreatedAt        string   `json:"created_at"`
	Distance         float64  `json:"distance,omitempty"`
}

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type MerchantItem struct {
	ItemId          string  `json:"item_id"`
	MerchantId      string  `json:"merchant_id,omitempty"`
	Name            string  `json:"name"`
	ProductCategory string  `json:"product_categories"`
	Price           float64 `json:"price"`
	ImageUrl        string  `json:"image_url"`
	CreatedAt       string  `json:"created_at"`
}
