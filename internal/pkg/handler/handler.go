package handler

import (
	"beli-mang/internal/db"
	"beli-mang/internal/pkg/configuration"
	"beli-mang/internal/pkg/service"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Run(cfg *configuration.Configuration, log *logrus.Logger) error {
	db, err := db.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	// set db to gin context
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	service := service.NewService(cfg, db)
	// userHandler := NewUserHandler(service)
	merchantHandler := NewMerchantHandler(service)
	imageHandler := NewImageHandler(service)

	// login
	// authGroup := router.Group("/v1/user/")
	// authGroup.POST("it/register", userHandler.Register)
	// authGroup.POST("it/login", userHandler.Login)
	// authGroup.POST("nurse/login", userHandler.Login)

	merchantGroup := router.Group("/admin/merchants/")
	// merchantGroup.Use(middleware.JWTAuth(cfg.JWTSecret, "admin"))
	merchantGroup.POST("", merchantHandler.AddMerchant)
	merchantGroup.POST(":merchantId/items", merchantHandler.AddMerchantItem)
	merchantGroup.GET("", merchantHandler.GetMerchants)

	imageGroup := router.Group("/image/")
	// imageGroup.Use(middleware.JWTAuth(cfg.JWTSecret, "admin"))
	imageGroup.POST("", imageHandler.UploadImage)

	return router.Run(":8080")
}
