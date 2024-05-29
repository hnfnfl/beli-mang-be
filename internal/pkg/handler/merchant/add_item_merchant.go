package merchant

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *MerchantHandler) AddMerchantItem(ctx *gin.Context) {
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

	item, errs := h.service.InsertMerchantItem(ctx, body)
	if errs.Code != 0 {
		errs.Send(ctx)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		item,
	)
}
