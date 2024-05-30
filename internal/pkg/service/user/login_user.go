package user

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) LoginUser(ctx *gin.Context, body model.User) *dto.AuthResponse {
	var (
		err error
		out model.User
	)

	db := s.db

	// check NIP in database
	stmt := "SELECT username, email, password_hash, role, email_role FROM users WHERE username = $1"
	if err := db.QueryRow(ctx, stmt, body.Username).Scan(
		&out.Username,
		&out.Email,
		&out.PasswordHash,
		&out.Role,
		&out.EmailRole,
	); err != nil {
		errs.NewNotFoundError(ctx, errs.ErrUserNotFound)
		return nil
	}

	//compare password
	if err := bcrypt.CompareHashAndPassword(out.PasswordHash, body.PasswordHash); err != nil {
		errs.NewBadRequestError(ctx, "password is wrong", err)
		return nil
	}

	// generate token
	var token string
	token, err = middleware.JWTSign(s.cfg, out.Username, out.Role)
	if err != nil {
		errs.NewInternalError(ctx, "token signing error", err)
		return nil
	}

	return &dto.AuthResponse{Token: token}
}
