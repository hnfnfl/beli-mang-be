package dto

import (
	"beli-mang/internal/db/model"
	"errors"
	"strconv"
	"sync"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	OrdersItems struct {
		MerchantId string
		ItemId     string
		Quantity   int
	}

	Orders struct {
		MerchantId      string        `json:"merchantId"`
		IsStartingPoint bool          `json:"isStartingPoint"`
		Items           []OrdersItems `json:"items"`
	}

	OrderEstimateRequest struct {
		UserLocation Location `json:"userLocation"`
		Orders       []Orders `json:"orders"`
	}

	OrderEstimateResponse struct {
		TotalPrice           float64 `json:"totalPrice"`
		EstDelivTime         int     `json:"estimatedDeliveryTimeInMinutes"`
		CalculatedEstimateId string  `json:"calculatedEstimateId"`
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

	GetOrdersRequest struct {
		MerchantId       string `form:"merchantId"`
		Limit            int    `form:"limit"`
		Offset           int    `form:"offset"`
		Name             string `form:"name"`
		MerchantCategory string `form:"merchantCategory"`
	}

	// Item struct
	OrderItemResponse struct {
		ItemId          string  `json:"itemId"`
		Name            string  `json:"name"`
		ProductCategory string  `json:"productCategory"`
		Price           float64 `json:"price"`
		Quantity        int     `json:"quantity"`
		ImageURL        string  `json:"imageUrl"`
		CreatedAt       string  `json:"createdAt"`
	}

	// Detail struct
	OrderResponse struct {
		Merchant model.Merchant      `json:"merchant"`
		Items    []OrderItemResponse `json:"items"`
	}

	GetOrdersResponse struct {
		OrderId string          `json:"orderId"`
		Orders  []OrderResponse `json:"orders"`
	}
)

func (r *OrderEstimateRequest) Validate() error {
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

	for _, order := range r.Orders {
		if len(order.Items) == 0 {
			return errors.New("itemId not found")
		}
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.Orders,
			validation.Required,
			validation.Each(validation.By(
				func(value interface{}) error {
					orders, ok := value.(Orders)
					if !ok {
						return validation.NewError("validation_OrderEstimateRequest", "invalid orders")
					}
					return orders.Validate()
				},
			)),
			validation.By(validateIsStartingPoint),
		),
	); err != nil {
		return err
	}

	return nil
}

func (r *Orders) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.Items,
			validation.Required,
			validation.Each(validation.By(
				func(value interface{}) error {
					ordersItems, ok := value.(OrdersItems)
					if !ok {
						return validation.NewError("validation_Orders", "invalid orders items")
					}
					return ordersItems.Validate()
				},
			)),
		),
		validation.Field(&r.MerchantId, validation.Required),
	); err != nil {
		return err
	}

	return nil
}

func (r *OrdersItems) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ItemId, validation.Required),
		validation.Field(&r.Quantity, validation.Required),
	); err != nil {
		return err
	}

	return nil
}

func validateIsStartingPoint(value interface{}) error {
	req, ok := value.([]Orders)
	if !ok {
		return validation.NewError("validation_orderEstimateRequest", "invalid order request")
	}

	count := 0
	for _, order := range req {
		if order.IsStartingPoint {
			count++
		}
	}
	if count != 1 {
		return validation.NewError("validation_isStartingPoint", "there must be exactly one starting point in orders")
	}

	return nil
}

func (r *PostOrderRequest) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.CalculatedEstimateId, validation.Required),
	); err != nil {
		return err
	}

	return nil
}

func (r *GetOrdersRequest) Validate() error {
	_, err := strconv.Atoi(r.MerchantId)
	if err == nil {
		return validation.NewError("merchantId", "merchantId must be a string")
	}

	_, err = strconv.Atoi(r.Name)
	if err == nil {
		return validation.NewError("name", "name must be a string")
	}

	if err := validation.Validate(&r.MerchantCategory,
		validation.In(MerchantCategoryList...),
	); err != nil {
		r.MerchantCategory = "<invalid>"
	}

	return nil

}
