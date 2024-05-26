package handler

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/service"
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageHandler struct {
	*service.Service
}

func NewImageHandler(s *service.Service) *ImageHandler {
	return &ImageHandler{s}
}

// UploadImage to s3 bucket
func (h *ImageHandler) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		errs.NewBadRequestError("file not found", err).Send(ctx)
		return
	}

	body := &dto.ImageRequest{
		File: file,
	}

	if err := body.Validate(); err != nil {
		errs.NewValidationError("Request validation error", err).Send(ctx)
		return
	}

	// Rename the file with UUID
	uuid := uuid.New()
	file.Filename = fmt.Sprintf("%s%s", uuid, filepath.Ext(file.Filename))

	// If the file passes all checks, you can continue with your processing
	h.UploadImageProcess(ctx, file).Send(ctx)
}
