package user

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/middleware"
	"beli-mang/internal/pkg/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Register(ctx *gin.Context) {
	body := dto.RegisterRequest{}
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

	passHash, err := middleware.PasswordHash(body.Password, h.handler.Config().Salt)
	if err != nil {
		errs.NewInternalError(ctx, "hashing error", err)
		return
	}

	data := model.User{
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: passHash,
	}

	role := util.ExtractRole(ctx.FullPath())
	switch role {
	case "admin":
		data.Role = "admin"
	case "users":
		data.Role = "user"
	}
	data.EmailRole = fmt.Sprintf("%s_%s", data.Email, data.Role)

	if token := h.service.RegisterUser(ctx, data); token != nil {
		ctx.JSON(http.StatusCreated, token)
	}
}
