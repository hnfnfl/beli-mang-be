package order

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) NearbyMerchant(ctx *gin.Context) {
	body := dto.GetNearbyMerchantsRequest{}
	msg, err := util.QueryBinding(ctx, &body)
	if err != nil {
		errs.NewValidationError(msg, err).Send(ctx)
		return
	}

	latlong := ctx.Param("latlong")
	latlongArr := strings.Split(latlong, ",")
	if len(latlongArr) != 2 {
		errs.NewValidationError("latlong", errs.ErrInvalidCoordinate).Send(ctx)
		return
	}

	// Convert latlongArr[0] to float64
	lat, err := strconv.ParseFloat(latlongArr[0], 64)
	if err != nil {
		errs.NewValidationError("lat", err).Send(ctx)
		return
	}
	body.Lat = lat

	long, err := strconv.ParseFloat(latlongArr[1], 64)
	if err != nil {
		errs.NewValidationError("long", err).Send(ctx)
		return
	}
	body.Long = long

	if body.Limit == 0 {
		body.Limit = 5
	}

	if err := body.Validate(); err != nil {
		errs.NewValidationError("Request validation error", err).Send(ctx)
		return
	}

	data, errs := h.service.GetNearbyMerchants(ctx, body)
	if errs.Code != 0 {
		errs.Send(ctx)
		return
	}

	ctx.JSON(
		http.StatusOK,
		data,
	)
}
