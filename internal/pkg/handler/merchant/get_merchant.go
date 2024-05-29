package merchant

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *MerchantHandler) GetMerchants(ctx *gin.Context) {
	body := dto.GetMerchantsRequest{}
	msg, err := util.QueryBinding(ctx, &body)
	if err != nil {
		errs.NewValidationError(msg, err).Send(ctx)
		return
	}

	body.CreatedAt = strings.ToUpper(body.CreatedAt)
	body.Name = strings.ToLower(body.Name)

	if body.Limit == 0 {
		body.Limit = 5
	}

	if err := body.Validate(); err != nil {
		errs.NewValidationError("Request validation error", err).Send(ctx)
		return
	}

	merchants, errs := h.service.GetMerchants(ctx, body)
	if errs.Code != 0 {
		errs.Send(ctx)
		return
	}

	ctx.JSON(
		http.StatusOK,
		merchants,
	)
}
