package middleware

import (
	"beli-mang/internal/pkg/configuration"
	"beli-mang/internal/pkg/errs"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func JWTSign(cfg *configuration.Configuration, username, role string) (string, error) {
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
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := ""

		if authHeader != "" {
			tokenString = authHeader[len(BearerSchema):]
		}

		if tokenString == "" {
			errs.NewUnauthorizedError("Authorization header not provided").Send(c)
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errs.ErrInvalidSigningMethod
			}

			return []byte(secret), nil
		})
		if err != nil {
			errs.NewUnauthorizedError(err.Error()).Send(c)
			c.Abort()
			return
		}

		if !token.Valid {
			errs.NewUnauthorizedError(errs.ErrInvalidToken.Error()).Send(c)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || claims.Issuer != expectedIssuer {
			errs.NewUnauthorizedError("Invalid User Role").Send(c)
			c.Abort()
			return
		}

		if claims.ExpiresAt.Before(time.Now()) {
			errs.NewUnauthorizedError(errs.ErrTokenExpired.Error()).Send(c)
			c.Abort()
			return
		}

		// c.Set("userID", claims.ID)
		// c.Set("userNIP", claims.Issuer)
		// c.Set("userRole", claims.Subject)

		c.Next()
	}
}

func PasswordHash(password string, salt string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
}
