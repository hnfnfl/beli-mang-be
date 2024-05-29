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
		errs.NewValidationError(msg, err).Send(ctx)
		return
	}

	// validate Request
	if err := body.Validate(); err != nil {
		errs.NewValidationError("Request validation error", err).Send(ctx)
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

	token, errs := h.service.LoginUser(ctx, data, *h.handler.Config())
	if errs.Code != 0 {
		errs.Send(ctx)
		return
	}

	ctx.JSON(http.StatusOK, token)
}
