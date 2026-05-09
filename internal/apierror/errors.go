package apierror

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

var (
	ErrUnauthorized        = errors.New("unauthorized")
	ErrValidation          = errors.New("validation error")
	ErrUsernameTaken       = errors.New("username already exists")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidToken        = errors.New("invalid token")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)

func ClassifyHTTPStatus(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, ErrUnauthorized), errors.Is(err, ErrInvalidToken):
		return http.StatusUnauthorized
	case errors.Is(err, ErrValidation), errors.Is(err, ErrInvalidPassword):
		return http.StatusBadRequest
	case errors.Is(err, ErrUsernameTaken):
		return http.StatusConflict
	case errors.Is(err, gorm.ErrRecordNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
