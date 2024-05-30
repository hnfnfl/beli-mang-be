package image

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (s *ImageService) UploadImage(ctx context.Context, file *multipart.FileHeader) errs.Response {
	s3Config := s.cfg.S3Config

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s3Config.Region),
		Credentials: credentials.NewStaticCredentials(
			s3Config.ID,
			s3Config.Secret,
			"",
		),
	})
	if err != nil {
		return errs.NewInternalError("failed to create aws session", err)
	}

	fileContent, err := file.Open()
	if err != nil {
		return errs.NewInternalError("failed to open file", err)
	}
	defer fileContent.Close()

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Body:   fileContent,
		Bucket: aws.String(s3Config.Bucket),
		Key:    aws.String(file.Filename),
	})
	if err != nil {
		return errs.NewInternalError("failed to upload image", err)
	}

	// return the image URL
	imageURL := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", s3Config.Bucket, s3Config.Region, file.Filename)

	return errs.Response{
		Code:    http.StatusOK,
		Message: "File uploaded successfully",
		Data: dto.ImageResponse{
			ImageUrl: imageURL,
		},
	}
}
