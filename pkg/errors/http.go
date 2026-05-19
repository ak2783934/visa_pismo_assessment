package apperrors

import (
	"errors"
	"net/http"
)

// HTTPStatus maps a custom (or wrapped) error to an HTTP status code.
func HTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrDuplicate):
		return http.StatusConflict
	case errors.Is(err, ErrValidation), errors.Is(err, ErrInvalidOperationType):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
