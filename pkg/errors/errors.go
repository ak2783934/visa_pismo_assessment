package apperrors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound             = errors.New("not found")
	ErrDuplicate            = errors.New("duplicate")
	ErrValidation           = errors.New("validation")
	ErrInvalidOperationType = errors.New("invalid operation type")
)

// Wrap returns an error that unwraps to kind and supports errors.Is against kind.
func Wrap(kind error, message string) error {
	if message == "" {
		return kind
	}
	return fmt.Errorf("%s: %w", message, kind)
}
