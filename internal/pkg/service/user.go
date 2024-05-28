package service

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/middleware"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) RegisterUser(body model.User) errs.Response {
	var err error

	db := s.DB()

	// insert user by role
	stmt := "INSERT INTO users (username, email, password_hash, role, email_role) VALUES ($1, $2, $3, $4, $5)"
	if _, err = db.Exec(stmt, body.Username, body.Email, body.PasswordHash, body.Role, body.EmailRole); err != nil {
		return errs.NewInternalError("insert error", err)
	}

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				errs.NewGenericError(http.StatusConflict, "user/admin is conflict")
			}
		}
		return errs.NewInternalError("insert error", err)
	}

	// generate token
	var token string
	if body.Role == "it" {
		token, err = middleware.JWTSign(s.Config(), body.Username, body.Role)
		if err != nil {
			return errs.NewInternalError("token signing error", err)
		}
	}

	// TODO: Implement this method
	return errs.Response{
		Code:    http.StatusOK,
		Message: "User registered successfully",
		Data: dto.AuthResponse{
			Token: token,
		},
	}
}

func (s *Service) LoginUser(body model.User) errs.Response {
	var (
		err error
		out model.User
	)

	db := s.DB()

	// check NIP in database
	stmt := "SELECT username, email, password_hash, role, email_role FROM users WHERE username = $1"
	if err := db.QueryRow(stmt, body.Username).Scan(
		&out.Username,
		&out.Email,
		&out.PasswordHash,
		&out.Role,
		&out.EmailRole,
	); err != nil {
		return errs.NewNotFoundError(errs.ErrUserNotFound)
	}

	//compare password
	if err := bcrypt.CompareHashAndPassword(out.PasswordHash, body.PasswordHash); err != nil {
		return errs.NewBadRequestError("password is wrong", err)
	}

	// generate token
	var token string
	token, err = middleware.JWTSign(s.Config(), out.Username, out.Role)
	if err != nil {
		return errs.NewInternalError("token signing error", err)
	}

	return errs.Response{
		Code:    http.StatusOK,
		Message: "User login successfully",
		Data: dto.AuthResponse{
			Token: token,
		},
	}
}
