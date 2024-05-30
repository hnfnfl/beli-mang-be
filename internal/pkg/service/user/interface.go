package user

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/configuration"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db  *pgxpool.Pool
	cfg *configuration.Configuration
}

type UserServiceInterface interface {
	RegisterUser(ctx *gin.Context, body model.User) (*dto.AuthResponse, errs.Response)
	LoginUser(ctx *gin.Context, body model.User) (*dto.AuthResponse, errs.Response)
}

func NewUserService(db *pgxpool.Pool, cfg *configuration.Configuration) *UserService {
	return &UserService{db, cfg}
}

var (
	_ UserServiceInterface = &UserService{}
)
