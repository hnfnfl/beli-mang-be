package dto

import (
	"sync"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	OrderEstimateItemPrice struct {
		ItemId string `json:"itemId"`
		Price  int    `json:"price"`
	}

	OrderEstimateRequestItem struct {
		ItemId   string `json:"itemId"`
		Quantity int    `json:"quantity"`
	}

	OrderEstimateRequestMerchant struct {
		MerchantId      string                     `json:"merchantId"`
		IsStartingPoint bool                       `json:"isStartingPoint"`
		Items           []OrderEstimateRequestItem `json:"items"`
	}

	OrderEstimateRequest struct {
		UserLocation Location                       `json:"userLocation"`
		Orders       []OrderEstimateRequestMerchant `json:"orders"`
	}

	OrderEstimateResponse struct {
		TotalPrice                     float64 `json:"totalPrice"`
		EstimatedDeliveryTimeInMinutes float64 `json:"estimatedDeliveryTimeInMinutes"`
		CalculatedEstimateId           string  `json:"calculatedEstimateId"`
	}

	CacheItem struct {
		Request  OrderEstimateRequest
		Response OrderEstimateResponse
		CachedAt time.Time
	}

	Cache struct {
		sync.RWMutex
		Data map[string]CacheItem
	}

	PostOrderRequest struct {
		CalculatedEstimateId string `json:"calculatedEstimateId"`
	}

	PostOrderResponse struct {
		OrderId string `json:"orderId"`
	}
)

func (r OrderEstimateRequestItem) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ItemId, validation.Required),
		validation.Field(&r.Quantity, validation.Required),
	); err != nil {
		return err
	}

	return nil
}

func (r OrderEstimateRequestMerchant) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.Items,
			validation.Required,
			validation.Each(validation.By(func(value interface{}) error {
				orderEstimateRequestMerchant, _ := value.(OrderEstimateRequestMerchant)
				return orderEstimateRequestMerchant.Validate()
			})),
		),
		validation.Field(&r.MerchantId, validation.Required),
		validation.Field(&r.IsStartingPoint, validation.Required),
	); err != nil {
		return err
	}

	return nil
}

func (r *OrderEstimateRequest) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.Orders,
			validation.Required,
			validation.Each(validation.By(func(value interface{}) error {
				orderEstimateRequestMerchant, _ := value.(OrderEstimateRequestMerchant)
				return orderEstimateRequestMerchant.Validate()
			})),
			validation.By(validateIsStartingPoint),
		),
		validation.Field(&r.UserLocation, validation.Required),
	); err != nil {
		return err
	}

	if err := validation.ValidateStruct(&r.UserLocation,
		validation.Field(&r.UserLocation.Lat, validation.Required),
		validation.Field(&r.UserLocation.Long, validation.Required),
	); err != nil {
		return err
	}

	// validate lat and long
	if r.UserLocation.Lat < -90 || r.UserLocation.Lat > 90 {
		return validation.NewError("lat", "latitude must be between -90 and 90")
	}

	if r.UserLocation.Long < -180 || r.UserLocation.Long > 180 {
		return validation.NewError("long", "longitude must be between -180 and 180")
	}

	return nil
}

func validateIsStartingPoint(value interface{}) error {
	req, _ := value.(OrderEstimateRequest)
	count := 0
	for _, order := range req.Orders {
		if order.IsStartingPoint {
			count++
		}
	}
	if count != 1 {
		return validation.NewError("validation_isStartingPoint", "there must be exactly one starting point in orders")
	}
	return nil
}

func (r PostOrderRequest) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.CalculatedEstimateId, validation.Required),
	); err != nil {
		return err
	}

	return nil
}
