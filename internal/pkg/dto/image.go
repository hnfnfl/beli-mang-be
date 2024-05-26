package dto

import (
	"beli-mang/internal/pkg/errs"
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ImageRequest struct {
	File *multipart.FileHeader
}

type ImageResponse struct {
	ImageUrl string `json:"imageUrl"`
}

func (i *ImageRequest) Validate() error {
	if err := validation.ValidateStruct(i,
		validation.Field(&i.File, validation.Required),
	); err != nil {
		return err
	}

	if i.File.Size < 10*1024 || i.File.Size > 2*1024*1024 {
		return validation.NewError("file", errs.ErrInvalidFileSize.Error())
	}

	if ext := i.File.Filename[len(i.File.Filename)-4:]; ext != ".jpg" && ext != "jpeg" {
		return validation.NewError("file", errs.ErrInvalidFileType.Error())
	}

	return nil
}
