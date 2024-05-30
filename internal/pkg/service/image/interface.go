package image

import (
	"beli-mang/internal/pkg/configuration"
	"beli-mang/internal/pkg/errs"
	"context"
	"mime/multipart"
)

type ImageService struct {
	cfg *configuration.Configuration
}

type ImageServiceInterface interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader) errs.Response
}

func NewImageService(cfg *configuration.Configuration) *ImageService {
	return &ImageService{cfg}
}

var (
	_ ImageServiceInterface = &ImageService{}
)