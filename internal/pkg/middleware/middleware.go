package middleware

import (
	"beli-mang/internal/pkg/configuration"
	"beli-mang/internal/pkg/errs"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func JWTSign(cfg *configuration.Configuration, username string, role string) (string, error) {
	expiry := time.Duration(cfg.AuthExpiry) * time.Hour
	timeStamp := time.Now()
	expiryTime := timeStamp.Add(expiry)

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeStamp),
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			ID:        username,
			Issuer:    role,
		},
	)

	return token.SignedString([]byte(cfg.JWTSecret))
}

func JWTAuth(secret string, expectedIssuer string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := ctx.GetHeader("Authorization")
		tokenString := ""

		if authHeader != "" {
			tokenString = authHeader[len(BearerSchema):]
		}

		if tokenString == "" {
			errs.NewUnauthorizedError(ctx, "Authorization header not provided")
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{},
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errs.ErrInvalidSigningMethod
				}

				return []byte(secret), nil
			},
		)
		if err != nil || token == nil {
			errs.NewUnauthorizedError(ctx, err.Error())
			return
		}

		if !token.Valid {
			errs.NewUnauthorizedError(ctx, errs.ErrInvalidToken.Error())
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || claims.Issuer != expectedIssuer {
			errs.NewUnauthorizedError(ctx, "Invalid User Role")
			return
		}

		if claims.ExpiresAt.Before(time.Now()) {
			errs.NewUnauthorizedError(ctx, errs.ErrTokenExpired.Error())
			return
		}

		ctx.Set("username", claims.ID)

		ctx.Next()
	}
}

func PasswordHash(password string, salt string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
}

func LoggerMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("logger", log)
		c.Next()
	}
}
