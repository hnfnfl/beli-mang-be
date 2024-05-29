package merchant

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	merchant, errs := h.service.InsertMerchant(ctx, body)
	if errs.Code != 0 {
		errs.Send(ctx)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		merchant,
	)
}
