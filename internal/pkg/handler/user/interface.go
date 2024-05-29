package user

import (
	"beli-mang/internal/pkg/handler"
	"beli-mang/internal/pkg/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	handler *handler.Handler
	service *service.Service
}

type UserInterface interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

func NewHandler(e *gin.Engine, h *handler.Handler, s *service.Service) {
	handler := &UserHandler{h, s}
	addRoutes(e, handler)
}

func addRoutes(r *gin.Engine, h UserInterface) {
	group := r.Group("")
	group.POST("/admin/register", h.Register)
	group.POST("/admin/login", h.Login)
	group.POST("/users/register", h.Register)
	group.POST("/users/login", h.Login)
}

var (
	_ UserInterface = &UserHandler{}
)
