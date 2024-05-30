package image

import (
	"beli-mang/internal/pkg/configuration"
	"beli-mang/internal/pkg/errs"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type ImageService struct {
	cfg *configuration.Configuration
}

type ImageServiceInterface interface {
	UploadImage(ctx *gin.Context, file *multipart.FileHeader) errs.Response
}

func NewImageService(cfg *configuration.Configuration) *ImageService {
	return &ImageService{cfg}
}

var (
	_ ImageServiceInterface = &ImageService{}
)
