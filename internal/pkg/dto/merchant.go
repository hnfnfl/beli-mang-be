package dto

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	MerchantCategory string
	ProductCategory  string
)

const (
	SmallRestaurant       MerchantCategory = "SmallRestaurant"
	MediumRestaurant      MerchantCategory = "MediumRestaurant"
	LargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	BoothKiosk            MerchantCategory = "BoothKiosk"
	ConvenienceStore      MerchantCategory = "ConvenienceStore"

	Beverage  ProductCategory = "Beverage"
	Food      ProductCategory = "Food"
	Snack     ProductCategory = "Snack"
	Condiment ProductCategory = "Condiment"
	Additions ProductCategory = "Additions"
)

var MerchantCategoryList = []interface{}{
	string(SmallRestaurant),
	string(MediumRestaurant),
	string(LargeRestaurant),
	string(MerchandiseRestaurant),
	string(BoothKiosk),
	string(ConvenienceStore),
}

var ProductCategoryList = []interface{}{
	string(Beverage),
	string(Food),
	string(Snack),
	string(Condiment),
	string(Additions),
}

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type AddMerchantRequest struct {
	MerchantId       string
	Name             string   `json:"name"`
	MerchantCategory string   `json:"merchantCategory"`
	ImageUrl         string   `json:"imageUrl"`
	Location         Location `json:"location"`
}

type AddMerchantResponse struct {
	MerchantId string `json:"merchantId"`
}

type AddMerchantItemRequest struct {
	ItemId          string
	MerchantId      string
	Name            string `json:"name"`
	ProductCategory string `json:"productCategory"`
	Price           int    `json:"price"`
	ImageUrl        string `json:"imageUrl"`
}

type AddMerchantItemResponse struct {
	ItemId string `json:"itemId"`
}

type GetMerchantsRequest struct {
	MerchantId       string `form:"merchantId"`
	Limit            int    `form:"limit"`
	Offset           int    `form:"offset"`
	Name             string `form:"name"`
	MerchantCategory string `form:"merchantCategory"`
	CreatedAt        string `form:"createdAt"`
}

type GetMerchantsResponse struct {
	Data []model.Merchant `json:"data"`
	Meta *errs.Meta       `json:"meta,omitempty"`
}

type GetMerchantItemsRequest struct {
	ItemId          string `form:"itemId"`
	Limit           int    `form:"limit"`
	Offset          int    `form:"offset"`
	Name            string `form:"name"`
	ProductCategory string `form:"productCategory"`
	CreatedAt       string `form:"createdAt"`
}

type GetMerchantItemsResponse struct {
	Data []model.MerchantItem `json:"data"`
	Meta *errs.Meta           `json:"meta,omitempty"`
}

func (r *AddMerchantRequest) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.Name,
			validation.Required,
			validation.Length(2, 30),
		),
		validation.Field(&r.MerchantCategory,
			validation.Required,
			validation.In(MerchantCategoryList...),
		),
		validation.Field(&r.ImageUrl,
			validation.Required,
			validation.NewStringRule(
				util.IsValidUrl, "must be a valid image url (jpg/jpeg)",
			),
		),
		validation.Field(&r.Location, validation.Required),
	); err != nil {
		return err
	}

	if err := validation.ValidateStruct(&r.Location,
		validation.Field(&r.Location.Lat, validation.Required),
		validation.Field(&r.Location.Long, validation.Required),
	); err != nil {
		return err
	}

	// validate lat and long
	if r.Location.Lat < -90 || r.Location.Lat > 90 {
		return validation.NewError("lat", "latitude must be between -90 and 90")
	}

	if r.Location.Long < -180 || r.Location.Long > 180 {
		return validation.NewError("long", "longitude must be between -180 and 180")
	}

	return nil
}

func (r *GetMerchantsRequest) Validate() error {
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

	if err := validation.Validate(&r.CreatedAt,
		validation.In(string(ASC), string(DESC)),
	); err != nil {
		r.CreatedAt = ""
	}

	return nil

}

func (r *AddMerchantItemRequest) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.MerchantId,
			validation.Required,
		),
		validation.Field(&r.Name,
			validation.Required,
			validation.Length(2, 30),
		),
		validation.Field(&r.ProductCategory,
			validation.Required,
			validation.In(ProductCategoryList...),
		),
		validation.Field(&r.Price,
			validation.Required,
			validation.Min(1),
		),
		validation.Field(&r.ImageUrl,
			validation.Required,
			validation.NewStringRule(
				util.IsValidUrl, "must be a valid image url (jpg/jpeg)",
			),
		),
	); err != nil {
		return err
	}

	return nil
}

func (r *GetMerchantItemsRequest) Validate() error {
	_, err := strconv.Atoi(r.ItemId)
	if err == nil {
		return validation.NewError("itemId", "itemId must be a string")
	}

	_, err = strconv.Atoi(r.Name)
	if err == nil {
		return validation.NewError("name", "name must be a string")
	}

	if err := validation.Validate(&r.ProductCategory,
		validation.In(ProductCategoryList...),
	); err != nil {
		r.ProductCategory = "<invalid>"
	}

	if err := validation.Validate(&r.CreatedAt,
		validation.In(string(ASC), string(DESC)),
	); err != nil {
		r.CreatedAt = ""
	}

	return nil
}
