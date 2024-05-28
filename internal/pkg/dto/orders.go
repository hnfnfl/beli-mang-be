package dto

type (
	ReqParamMerchantNearbyGet struct {
		Name             string           `json:"name"`
		MerchantId       string           `json:"merchantId"`
		MerchantCategory MerchantCategory `json:"merchantCategory"`
		Limit            int              `json:"limit"`
		Offset           int              `json:"offset"`
	}
	ResMerchantNearbyMerchantsGet struct {
		MerchantId       string           `json:"merchantId"`
		Name             string           `json:"name"`
		ImageUrl         string           `json:"imageUrl"`
		Location         Location         `json:"location"`
		MerchantCategory MerchantCategory `json:"merchantCategory"`
		CreatedAt        string           `json:"createdAt"`
	}
	ResMerchantNearbyItemsGet struct {
		ItemId          string           `json:"itemId"`
		Name            string           `json:"name"`
		ProductCategory MerchantCategory `json:"productCategory"`
		Price           string           `json:"price"`
		ImageUrl        string           `json:"imageUrl"`
		CreatedAt       string           `json:"createdAt"`
	}
	ResMerchantNearbyGet struct {
		Merchant ResMerchantNearbyMerchantsGet `json:"merchant"`
		Items    []ResMerchantNearbyItemsGet   `json:"items"`
	}
)
