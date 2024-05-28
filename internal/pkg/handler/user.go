package handler

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/middleware"
	"beli-mang/internal/pkg/service"
	"beli-mang/internal/pkg/util"
	"fmt"
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

	data := model.User{
		Username: body.Username,
		Email:    body.Email,
	}

	var passHash []byte
	if body.Password != "" {
		var err error
		passHash, err = middleware.PasswordHash(body.Password, h.service.Config().Salt)
		if err != nil {
			errs.NewInternalError("hashing error", err).Send(ctx)
			return
		}
	}

	role := extractRole(ctx.FullPath())
	data.PasswordHash = passHash

	switch role {
	case "admin":
		data.Role = "admin"
	case "users":
		data.Role = "user"
	}
	data.EmailRole = fmt.Sprintf("%s_%s", data.Email, data.Role)
	h.service.RegisterUser(data).Send(ctx)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	body := dto.LoginRequest{}
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

	data := model.User{
		Username: body.Username,
	}

	if body.Password != "" {
		data.PasswordHash = []byte(body.Password + h.service.Config().Salt)
	}

	role := extractRole(ctx.FullPath())

	switch role {
	case "admin":
		data.Role = "admin"
		h.service.RegisterUser(data).Send(ctx)
	case "users":
		data.Role = "user"
		h.service.RegisterUser(data).Send(ctx)
	}
}

func extractRole(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}
