package handler

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/service"
	"beli-mang/internal/pkg/util"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.Service
}

func NewUserHandler(s *service.Service) *UserHandler {
	return &UserHandler{s}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	body := dto.RegisterRequest{}
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

	// data := model.User{
	// }

	// var passHash []byte
	// if body.Password != "" {
	// 	var err error
	// 	passHash, err = middleware.PasswordHash(body.Password, h.service.Config().Salt)
	// 	if err != nil {
	// 		errs.NewInternalError("hashing error", err).Send(ctx)
	// 		return
	// 	}
	// }
}

func (h *UserHandler) Login(ctx *gin.Context) {
	body := dto.LoginRequest{}
	msg, err := util.JsonBinding(ctx, &body)
	if err != nil {
		errs.NewValidationError(msg, err).Send(ctx)
		return
	}

	// h.service.LoginUser(data).Send(ctx)
}

func extractRole(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 4 {
		return parts[3]
	}
	return ""
}
