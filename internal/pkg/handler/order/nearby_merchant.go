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
		errs.NewValidationError(ctx, msg, err)
		return
	}

	latlong := ctx.Param("latlong")
	latlongArr := strings.Split(latlong, ",")
	if len(latlongArr) != 2 {
		errs.NewValidationError(ctx, "latlong", errs.ErrInvalidCoordinate)
		return
	}

	// Convert latlongArr[0] to float64
	lat, err := strconv.ParseFloat(latlongArr[0], 64)
	if err != nil {
		errs.NewValidationError(ctx, "lat", err)
		return
	}
	body.Lat = lat

	long, err := strconv.ParseFloat(latlongArr[1], 64)
	if err != nil {
		errs.NewValidationError(ctx, "long", err)
		return
	}
	body.Long = long

	if body.Limit == 0 {
		body.Limit = 5
	}

	if err := body.Validate(); err != nil {
		errs.NewValidationError(ctx, "Request validation error", err)
		return
	}

	if data := h.service.GetNearbyMerchants(ctx, body); data != nil {
		ctx.JSON(http.StatusOK, data)
	}
}
