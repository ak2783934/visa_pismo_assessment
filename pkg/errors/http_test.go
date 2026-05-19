package apperrors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPStatus(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
	}{
		{name: "not found", err: Wrap(ErrNotFound, "account not found"), status: http.StatusNotFound},
		{name: "duplicate", err: Wrap(ErrDuplicate, "already exists"), status: http.StatusConflict},
		{name: "validation", err: Wrap(ErrValidation, "required"), status: http.StatusBadRequest},
		{name: "invalid operation", err: ErrInvalidOperationType, status: http.StatusBadRequest},
		{name: "unknown", err: assert.AnError, status: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.status, HTTPStatus(tt.err))
		})
	}
}
