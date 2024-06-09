package order

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) EstimateOrder(ctx *gin.Context) {
	body := dto.OrderEstimateRequest{}
	msg, err := util.JsonBinding(ctx, &body)
	if err != nil {
		errs.NewValidationError(ctx, msg, err)
		return
	}

	// validate Request
	if err := body.Validate(); err != nil {
		if err.Error() == "itemId not found" {
			errs.NewNotFoundError(ctx, errs.ErrItemNotFound)
			return
		}

		errs.NewValidationError(ctx, "Request validation error", err)
		return
	}

	if res := h.service.EstimateOrder(ctx, body); res != nil {
		ctx.JSON(http.StatusOK, res)
	}
}
