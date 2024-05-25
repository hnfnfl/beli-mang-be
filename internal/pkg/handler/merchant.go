package handler

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/service"
	"beli-mang/internal/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	service *service.Service
}

func NewMerchantHandler(s *service.Service) *MerchantHandler {
	return &MerchantHandler{s}
}

func (h *MerchantHandler) AddMerchant(ctx *gin.Context) {
	body := dto.AddMerchantRequest{}
	msg, err := util.JsonBinding(ctx, &body)
	if err != nil {
		errs.NewValidationError(msg, err).Send(ctx)
		return
	}

	// validate Request
	if err := body.Validate(); err != nil {
		errs.NewValidationError("Request validation error", err).Send(ctx)
		return
	}

	var prefixID string
	switch body.MerchantCategory {
	case string(dto.SmallRestaurant):
		prefixID = "SR"
	case string(dto.MediumRestaurant):
		prefixID = "MR"
	case string(dto.LargeRestaurant):
		prefixID = "LR"
	case string(dto.MerchandiseRestaurant):
		prefixID = "ME"
	case string(dto.BoothKiosk):
		prefixID = "BK"
	case string(dto.ConvenienceStore):
		prefixID = "CS"
	default:
		prefixID = "MC"
	}

	body.MerchantId = util.UuidGenerator(prefixID, 15)

	merchantID, errs := h.service.InsertMerchant(body)
	if errs.Code != 0 {
		errs.Send(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, dto.AddMerchantResponse{MerchantId: merchantID})
}

func (h *MerchantHandler) AddMerchantItems(ctx *gin.Context) {
	body := dto.AddMerchantItemRequest{}
	msg, err := util.JsonBinding(ctx, &body)
	if err != nil {
		errs.NewValidationError(msg, err).Send(ctx)
		return
	}

	body.MerchantId = ctx.Param("merchantId")
	
	var prefixID string
	switch body.ProductCategory {
	case string(dto.Beverage):
		prefixID = "B"
	case string(dto.Food):
		prefixID = "F"
	case string(dto.Snack):
		prefixID = "S"
	case string(dto.Condiment):
		prefixID = "C"
	case string(dto.Additions):
		prefixID = "A"
	default:
		prefixID = "I"
	}

	prefixID += body.MerchantId[:2]
	body.ItemId = util.UuidGenerator(prefixID, 15)

	// validate Request
	if err := body.Validate(); err != nil {
		errs.NewValidationError("Request validation error", err).Send(ctx)
		return
	}

	merchantItemID, errs := h.service.InsertMerchantItems(body)
	if errs.Code != 0 {
		errs.Send(ctx)
		return
	}

	ctx.JSON(200, dto.AddMerchantItemsResponse{ItemId: merchantItemID})
}
