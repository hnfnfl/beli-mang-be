package server

import (
	"beli-mang/internal/db"
	"beli-mang/internal/pkg/configuration"
	"beli-mang/internal/pkg/handler"
	"beli-mang/internal/pkg/handler/image"
	"beli-mang/internal/pkg/handler/merchant"
	"beli-mang/internal/pkg/handler/order"
	"beli-mang/internal/pkg/handler/user"
	"beli-mang/internal/pkg/service"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Run(cfg *configuration.Configuration, log *logrus.Logger) error {
	ctx := context.Background()
	db := db.GetConn(cfg, ctx)

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	service := service.NewService(db)
	handler := handler.NewHandler(cfg, log)

	user.NewHandler(router, handler, service)
	merchant.NewHandler(router, handler, service)
	image.NewHandler(router, handler, service)
	order.NewHandler(router, handler, service)

	return router.Run(":8080")
}
