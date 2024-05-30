package order

import (
	"beli-mang/internal/pkg/handler"
	"beli-mang/internal/pkg/middleware"
	"beli-mang/internal/pkg/service/order"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	handler *handler.Handler
	service *order.OrderService
}

type OrderInterface interface {
	NearbyMerchant(ctx *gin.Context)
	EstimateOrder(ctx *gin.Context)
	// AddOrders(ctx *gin.Context)
	// GetOrders(ctx *gin.Context)
}

func NewHandler(e *gin.Engine, h *handler.Handler) {
	s := order.NewOrderService(h.DB())
	handler := &OrderHandler{h, s}
	secret := h.Config().JWTSecret
	addRoutes(e, handler, secret)
}

func addRoutes(r *gin.Engine, h OrderInterface, secret string) {
	group := r.Group("")
	group.Use(middleware.JWTAuth(secret, "user"))
	group.GET("/merchants/nearby/:latlong", h.NearbyMerchant)
	group.POST("/users/estimate", h.EstimateOrder)
	// group.POST("/users/orders", h.AddOrders)
	// group.GET("/users/orders", h.GetOrders)
}

var (
	_ OrderInterface = &OrderHandler{}
)
