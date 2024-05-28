package service

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/middleware"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) RegisterUser(body model.User) (*dto.AuthResponse, errs.Response) {
	var err error
	db := s.DB()

	// check if username and email with same role is already exist
	var usernameCount, emailCount int
	stmt := `SELECT 
    (SELECT COUNT(*) FROM users WHERE username = $1) AS username_count,
    (SELECT COUNT(*) FROM users WHERE email = $2 AND role = $3) AS email_count`
	if err := db.QueryRow(stmt, body.Username, body.Email, body.Role).Scan(&usernameCount, &emailCount); err != nil {
		return nil, errs.NewInternalError("check user error", err)
	}

	if usernameCount > 0 {
		return nil, errs.NewGenericError(http.StatusConflict, "username is already exist")
	}

	if emailCount > 0 {
		return nil, errs.NewGenericError(http.StatusConflict, "email is already exist in the same role")
	}

	// insert user by role
	stmt = "INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, $4)"
	if _, err := db.Exec(stmt, body.Username, body.Email, body.PasswordHash, body.Role); err != nil {
		return nil, errs.NewInternalError("insert error", err)
	}

	token, err := middleware.JWTSign(s.Config(), body.Username, body.Role)
	if err != nil {
		return nil, errs.NewInternalError("token signing error", err)
	}

	return &dto.AuthResponse{Token: token}, errs.Response{}
}

func (s *Service) LoginUser(body model.User) errs.Response {
	var (
		err error
		out model.User
	)

	db := s.DB()

	// check NIP in database
	stmt := "SELECT username, email, password_hash, role FROM users WHERE username = $1"
	if err := db.QueryRow(stmt, body.Username).Scan(
		&out.Username,
		&out.Email,
		&out.PasswordHash,
		&out.Role,
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
