package image

import (
	"beli-mang/internal/pkg/handler"
	"beli-mang/internal/pkg/middleware"
	"beli-mang/internal/pkg/service/image"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	handler *handler.Handler
	service *image.ImageService
}

type ImageInterface interface {
	UploadImage(ctx *gin.Context)
}

func NewHandler(e *gin.Engine, h *handler.Handler) {
	s := image.NewImageService(h.Config())
	handler := &ImageHandler{h, s}
	secret := h.Config().JWTSecret
	addRoutes(e, handler, secret)
}

func addRoutes(r *gin.Engine, h ImageInterface, secret string) {
	group := r.Group("/image")
	group.Use(middleware.JWTAuth(secret, "admin"))
	group.POST("", h.UploadImage)
}

var (
	_ ImageInterface = &ImageHandler{}
)
