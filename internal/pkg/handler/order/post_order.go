package order

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) PostOrder(ctx *gin.Context) {
	body := dto.PostOrderRequest{}
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

	if res := h.service.PostOrder(ctx, body); res != nil {
		ctx.JSON(http.StatusCreated, res)
	}
}
