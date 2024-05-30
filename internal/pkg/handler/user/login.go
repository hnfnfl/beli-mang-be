package user

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Login(ctx *gin.Context) {
	body := dto.LoginRequest{}
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

	data := model.User{
		Username:     body.Username,
		PasswordHash: []byte(body.Password + h.handler.Config().Salt),
	}
	role := util.ExtractRole(ctx.FullPath())

	switch role {
	case "admin":
		data.Role = "admin"
	case "users":
		data.Role = "user"
	}

	if token := h.service.LoginUser(ctx, data); token != nil {
		ctx.JSON(http.StatusOK, token)
	}
}
