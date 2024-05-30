package errs

import (
	"beli-mang/internal/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type Response struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Err     string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

func (e Response) Error() string {
	return e.Err
}

func finishWithError(ctx *gin.Context, resp Response) {
	ctx.Abort()
	_ = ctx.Error(resp)
}

func NewGenericError(ctx *gin.Context, code int, msg string) {
	finishWithError(ctx, Response{
		Code:    code,
		Message: msg,
	})
}

func NewInternalError(ctx *gin.Context, msg string, err error) {
	finishWithError(ctx, Response{
		Code:    http.StatusInternalServerError,
		Err:     err.Error(),
		Message: msg,
	})
}

func NewNotFoundError(ctx *gin.Context, err error) {
	finishWithError(ctx, Response{
		Code: http.StatusNotFound,
		Err:  err.Error(),
	})
}

func NewValidationError(ctx *gin.Context, msg string, err error) {
	finishWithError(ctx, Response{
		Code: http.StatusBadRequest,
		Err:  err.Error(),
	})
}

func NewBadRequestError(ctx *gin.Context, msg string, err error) {
	finishWithError(ctx, Response{
		Code:    http.StatusBadRequest,
		Err:     err.Error(),
		Message: msg,
	})
}

func NewUnauthorizedError(ctx *gin.Context, msg string) {
	finishWithError(ctx, Response{
		Code:    http.StatusUnauthorized,
		Message: msg,
	})
}

func (e Response) Send(ctx *gin.Context) {
	ctx.JSON(e.Code, e)
}

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors[0].Err
			switch err := err.(type) {
			case Response:
				err.Send(ctx)
			default:
				logger.FromContext(ctx).Error(err)
				ctx.JSON(http.StatusInternalServerError, err)
			}
		}
	}
}
