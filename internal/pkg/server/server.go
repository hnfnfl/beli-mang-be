package server

import (
	"beli-mang/internal/db"
	"beli-mang/internal/pkg/configuration"
	"beli-mang/internal/pkg/handler"
	"beli-mang/internal/pkg/handler/image"
	"beli-mang/internal/pkg/handler/merchant"
	"beli-mang/internal/pkg/handler/order"
	"beli-mang/internal/pkg/handler/user"
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

	handler := handler.NewHandler(cfg, db, log)

	user.NewHandler(router, handler)
	merchant.NewHandler(router, handler)
	image.NewHandler(router, handler)
	order.NewHandler(router, handler)

	return router.Run(":8080")
}
