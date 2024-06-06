package model

type Merchant struct {
	MerchantId       string   `json:"merchantId"`
	Name             string   `json:"name"`
	MerchantCategory string   `json:"merchantCategory"`
	ImageUrl         string   `json:"imageUrl"`
	Location         Location `json:"location"`
	CreatedAt        string   `json:"createdAt"`
	Distance         float64  `json:"distance,omitempty"`
}

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type MerchantItem struct {
	ItemId          string  `json:"itemId"`
	MerchantId      string  `json:"merchantId,omitempty"`
	Name            string  `json:"name"`
	ProductCategory string  `json:"productCategory"`
	Price           float64 `json:"price"`
	ImageUrl        string  `json:"imageUrl"`
	CreatedAt       string  `json:"createdAt"`
}
