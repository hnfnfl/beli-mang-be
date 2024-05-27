package errs

import (
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
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

func NewGenericError(code int, msg string) Response {
	return Response{
		Code:    code,
		Message: msg,
	}
}

func NewInternalError(msg string, err error) Response {
	return Response{
		Code:    http.StatusInternalServerError,
		Error:   err.Error(),
		Message: msg,
	}
}

func NewNotFoundError(err error) Response {
	return Response{
		Code:  http.StatusNotFound,
		Error: err.Error(),
	}
}

func NewValidationError(msg string, err error) Response {
	return Response{
		Code:    http.StatusBadRequest,
		Error:   err.Error(),
		Message: msg,
	}
}

func NewBadRequestError(msg string, err error) Response {
	return Response{
		Code:    http.StatusBadRequest,
		Error:   err.Error(),
		Message: msg,
	}
}

func NewUnauthorizedError(msg string) Response {
	return Response{
		Code:    http.StatusUnauthorized,
		Message: msg,
	}
}

func (e Response) Send(ctx *gin.Context) {
	ctx.JSON(e.Code, e)
}
