package errs

import "errors"

// nolint::lll
var (
	ErrUserNotFound         = errors.New("user not found")
	ErrMerchantNotFound     = errors.New("merchant not found")
	ErrPasswordMissmatch    = errors.New("password missmatch")
	ErrInvalidClaimsType    = errors.New("invalid claims type")
	ErrInvalidToken         = errors.New("invalid token")
	ErrUnauthorized         = errors.New("user is not authorized")
	ErrInvalidSigningMethod = errors.New("invalid signing method algorithm")
	ErrTokenExpired         = errors.New("token expired")
	ErrInvalidFileType      = errors.New("invalid file type")
	ErrInvalidFileSize      = errors.New("invalid file size")
	ErrBadParam             = errors.New("param request is invalid")
	ErrInternalServerError  = errors.New("internal server error")
	ErrInvalidCoordinate    = errors.New("invalid coordinate format")
	ErrMerchantTooFar       = errors.New("distance between user and merchant is more than 3km")
)
