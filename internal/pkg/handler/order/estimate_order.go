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
		errs.NewValidationError(msg, err).Send(ctx)
		return
	}

	// validate Request
	if err := body.Validate(); err != nil {
		errs.NewValidationError("Request validation error", err).Send(ctx)
		return
	}

	var (
		res  *dto.OrderEstimateResponse
		errs errs.Response
	)

	res, errs = h.service.EstimateOrder(ctx, body)
	if errs.Code != 0 {
		errs.Send(ctx)
		return
	}
	ctx.JSON(http.StatusOK, res)
}
