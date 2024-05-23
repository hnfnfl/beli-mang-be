package service

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"net/http"
)

func (s *Service) RegisterUser(body model.User) errs.Response {
	// TODO: Implement this method
	return errs.Response{
		Code:    http.StatusOK,
		Message: "User registered successfully",
		Data:    dto.AuthResponse{},
	}
}

func (s *Service) LoginUser(body model.User) errs.Response {
	// TODO: Implement this method

	return errs.Response{
		Code:    http.StatusOK,
		Message: "User login successfully",
		Data:    dto.AuthResponse{},
	}
}
