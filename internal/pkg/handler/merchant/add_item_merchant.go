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
		errs.NewValidationError(ctx, msg, err)
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
	case string(dto.Condiments):
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
		errs.NewValidationError(ctx, "Request validation error", err)
		return
	}

	item := h.service.InsertMerchantItem(ctx, body)
	if item != nil {
		ctx.JSON(http.StatusCreated, item)
	}
}
