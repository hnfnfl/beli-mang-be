package image

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
	h.service.UploadImage(ctx, file).Send(ctx)
}
