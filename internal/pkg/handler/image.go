package handler

import (
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
	// The "image" here is the name of the form field
	file, err := ctx.FormFile("file")
	if err != nil {
		errs.NewBadRequestError("file not found", err).Send(ctx)
		return
	}

	// Check the file size (no more than 2MB, no less than 10KB)
	if file.Size < 10*1024 || file.Size > 2*1024*1024 {
		errs.NewBadRequestError(
			fmt.Sprintf("file size: %d bytes", file.Size),
			errs.ErrInvalidFileSize,
		).Send(ctx)
		return
	}

	// Check the file extension (must be .jpg or .jpeg)
	if ext := filepath.Ext(file.Filename); ext != ".jpg" && ext != ".jpeg" {
		errs.NewBadRequestError(
			fmt.Sprintf("file extension: %s", ext),
			errs.ErrInvalidFileType,
		).Send(ctx)
		return
	}

	// Rename the file with UUID
	uuid := uuid.New()
	file.Filename = fmt.Sprintf("%s%s", uuid, filepath.Ext(file.Filename))

	// If the file passes all checks, you can continue with your processing
	h.UploadImageProcess(ctx, file).Send(ctx)
}
