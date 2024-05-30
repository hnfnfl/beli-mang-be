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
		errs.NewValidationError(ctx, msg, err)
		return
	}

	// validate Request
	if err := body.Validate(); err != nil {
		errs.NewValidationError(ctx, "Request validation error", err)
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

	if merchant := h.service.InsertMerchant(ctx, body); merchant != nil {
		ctx.JSON(http.StatusCreated, merchant)
	}
}
