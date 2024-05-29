package merchant

import (
	"beli-mang/internal/pkg/handler"
	"beli-mang/internal/pkg/middleware"
	"beli-mang/internal/pkg/service"

	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	handler *handler.Handler
	service *service.Service
}

type MerchantInterface interface {
	AddMerchant(ctx *gin.Context)
	GetMerchants(ctx *gin.Context)
	AddMerchantItem(ctx *gin.Context)
	GetMerchantItems(ctx *gin.Context)
}

func NewHandler(e *gin.Engine, h *handler.Handler, s *service.Service) {
	handler := &MerchantHandler{h, s}
	secret := h.Config().JWTSecret
	addRoutes(e, handler, secret)
}

func addRoutes(r *gin.Engine, h MerchantInterface, secret string) {
	group := r.Group("/admin/merchants")
	group.Use(middleware.JWTAuth(secret, "admin"))
	group.POST("", h.AddMerchant)
	group.GET("", h.GetMerchants)
	group.POST("/:merchantId/items", h.AddMerchantItem)
	group.GET("/:merchantId/items", h.GetMerchantItems)
}

var (
	_ MerchantInterface = &MerchantHandler{}
)
