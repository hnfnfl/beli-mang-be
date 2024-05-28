package handler

import (
	"beli-mang/internal/db"
	"beli-mang/internal/pkg/configuration"
	"beli-mang/internal/pkg/middleware"
	"beli-mang/internal/pkg/service"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Run(cfg *configuration.Configuration, log *logrus.Logger) error {
	ctx := context.Background()
	// conn := postgre.GetConn(ctx)
	db := db.GetConn(cfg, ctx)

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	// set db to gin context
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	service := service.NewService(cfg, db)
	userHandler := NewUserHandler(service)
	merchantHandler := NewMerchantHandler(service)
	imageHandler := NewImageHandler(service)
	// orderHandler := NewOrderHandler(service)

	// login
	authGroup := router.Group("")
	authGroup.POST("/admin/register", userHandler.Register)
	authGroup.POST("/admin/login", userHandler.Login)
	authGroup.POST("/users/register", userHandler.Register)
	authGroup.POST("/users/login", userHandler.Login)

	merchantGroup := router.Group("/admin/merchants/")
	merchantGroup.Use(middleware.JWTAuth(cfg.JWTSecret, "admin"))
	merchantGroup.POST("", merchantHandler.AddMerchant)
	merchantGroup.GET("", merchantHandler.GetMerchants)
	merchantGroup.POST(":merchantId/items", merchantHandler.AddMerchantItem)
	merchantGroup.GET(":merchantId/items", merchantHandler.GetMerchantItems)

	imageGroup := router.Group("/image/")
	// imageGroup.Use(middleware.JWTAuth(cfg.JWTSecret, "admin"))
	imageGroup.POST("", imageHandler.UploadImage)

	return router.Run(":8080")
}
