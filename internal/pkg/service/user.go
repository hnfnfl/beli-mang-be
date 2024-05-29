package service

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/middleware"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) RegisterUser(ctx *gin.Context, body model.User) (*dto.AuthResponse, errs.Response) {
	// var err error

	db := s.DB()

	// insert user by role
	// var userName string
	// var emailRole string
	stmt := "INSERT INTO users (username, email, password_hash, role, email_role) VALUES ($1, $2, $3, $4, $5) RETURNING username, email_role"
	_, err := db.Exec(ctx, stmt, body.Username, body.Email, body.PasswordHash, body.Role, body.EmailRole)
	fmt.Println(err)
	pgErr, ok := err.(*pgconn.PgError)
	fmt.Println(pgErr, ok)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				if pgErr.ConstraintName == "users_pkey" {
					return nil, errs.NewGenericError(http.StatusConflict, "user/admin username is conflict")
				} else if pgErr.ConstraintName == "users_email_role_key" {
					return nil, errs.NewGenericError(http.StatusConflict, "user/admin email is conflict")
				}
			}
		}
		return nil, errs.NewInternalError("insert error", err)
	}

	// generate token
	var token string
	token, err = middleware.JWTSign(s.Config(), body.Username, body.Role)
	if err != nil {
		return nil, errs.NewInternalError("token signing error", err)
	}

	return &dto.AuthResponse{Token: token}, errs.Response{}
}

func (s *Service) LoginUser(ctx *gin.Context, body model.User) (*dto.AuthResponse, errs.Response) {
	var (
		err error
		out model.User
	)

	db := s.DB()

	// check NIP in database
	stmt := "SELECT username, email, password_hash, role, email_role FROM users WHERE username = $1"
	if err := db.QueryRow(ctx, stmt, body.Username).Scan(
		&out.Username,
		&out.Email,
		&out.PasswordHash,
		&out.Role,
		&out.EmailRole,
	); err != nil {
		return nil, errs.NewNotFoundError(errs.ErrUserNotFound)
	}

	//compare password
	if err := bcrypt.CompareHashAndPassword(out.PasswordHash, body.PasswordHash); err != nil {
		return nil, errs.NewBadRequestError("password is wrong", err)
	}

	// generate token
	var token string
	token, err = middleware.JWTSign(s.Config(), out.Username, out.Role)
	if err != nil {
		return nil, errs.NewInternalError("token signing error", err)
	}

	return &dto.AuthResponse{Token: token}, errs.Response{}
}
