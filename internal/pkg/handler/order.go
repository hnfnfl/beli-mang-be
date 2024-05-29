package handler

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/service"
	"beli-mang/internal/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *service.Service
}

func NewOrderHandler(s *service.Service) *OrderHandler {
	return &OrderHandler{s}
}

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

// func (h *OrderHandler) Register(ctx *gin.Context) {
// 	body := dto.RegisterRequest{}
// 	msg, err := util.JsonBinding(ctx, &body)
// 	if err != nil {
// 		errs.NewValidationError(msg, err).Send(ctx)
// 		return
// 	}

// 	// validate Request
// 	if err := body.Validate(); err != nil {
// 		errs.NewValidationError("Request validation error", err).Send(ctx)
// 		return
// 	}

// 	data := model.Order{
// 		Ordername: body.Ordername,
// 		Email:     body.Email,
// 	}

// 	var passHash []byte
// 	if body.Password != "" {
// 		var err error
// 		passHash, err = middleware.PasswordHash(body.Password, h.service.Config().Salt)
// 		if err != nil {
// 			errs.NewInternalError("hashing error", err).Send(ctx)
// 			return
// 		}
// 	}

// 	role := extractRole(ctx.FullPath())
// 	data.PasswordHash = passHash

// 	switch role {
// 	case "admin":
// 		data.Role = "admin"
// 		h.service.RegisterOrder(data).Send(ctx)
// 	case "users":
// 		data.Role = "user"
// 		h.service.RegisterOrder(data).Send(ctx)
// 	}
// }

// func (h *OrderHandler) Login(ctx *gin.Context) {
// 	body := dto.LoginRequest{}
// 	msg, err := util.JsonBinding(ctx, &body)
// 	if err != nil {
// 		errs.NewValidationError(msg, err).Send(ctx)
// 		return
// 	}

// 	// validate Request
// 	if err := body.Validate(); err != nil {
// 		errs.NewValidationError("Request validation error", err).Send(ctx)
// 		return
// 	}

// 	data := model.Order{
// 		Ordername: body.Ordername,
// 	}

// 	if body.Password != "" {
// 		data.PasswordHash = []byte(body.Password + h.service.Config().Salt)
// 	}

// 	role := extractRole(ctx.FullPath())

// 	switch role {
// 	case "admin":
// 		data.Role = "admin"
// 		h.service.RegisterOrder(data).Send(ctx)
// 	case "users":
// 		data.Role = "user"
// 		h.service.RegisterOrder(data).Send(ctx)
// 	}
// }
