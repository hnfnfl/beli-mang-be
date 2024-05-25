package handler

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/service"
	"beli-mang/internal/pkg/util"

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

	data := model.Merchant{
		Name:             body.Name,
		MerchantCategory: string(body.MerchantCategory),
		ImageUrl:         body.ImageUrl,
		Lat:              body.Location.Lat,
		Long:             body.Location.Long,
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
		prefixID = "MeR"
	case string(dto.BoothKiosk):
		prefixID = "BK"
	case string(dto.ConvenienceStore):
		prefixID = "CS"
	default:
		prefixID = "M"
	}

	data.MerchantId = util.UuidGenerator(prefixID, 15)

	merchantID, errs := h.service.InsertMerchant(data)
	if errs.Code != 0 {
		errs.Send(ctx)
		return
	}

	ctx.JSON(200, gin.H{
		"merchantId": merchantID,
	})
}
