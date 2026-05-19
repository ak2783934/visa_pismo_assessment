package apperrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorsIs(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		err := Wrap(ErrNotFound, "account not found")
		assert.True(t, errors.Is(err, ErrNotFound))
		assert.False(t, errors.Is(err, ErrDuplicate))
	})

	t.Run("duplicate", func(t *testing.T) {
		err := Wrap(ErrDuplicate, "document number already exists")
		assert.True(t, errors.Is(err, ErrDuplicate))
	})

	t.Run("validation", func(t *testing.T) {
		err := Wrap(ErrValidation, "amount must be greater than 0")
		assert.True(t, errors.Is(err, ErrValidation))
	})

	t.Run("invalid operation type", func(t *testing.T) {
		assert.True(t, errors.Is(ErrInvalidOperationType, ErrInvalidOperationType))
	})
}
