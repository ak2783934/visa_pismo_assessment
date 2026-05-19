package response

import (
	"net/http"

	apperrors "github.com/ak2783934/visa-pismo-assessment/pkg/errors"
	"github.com/gin-gonic/gin"
)

func JSON(c *gin.Context, status int, payload any) {
	c.JSON(status, payload)
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// Handle writes a JSON error from err using apperrors.HTTPStatus.
// For unexpected errors (500), fallbackMessage is returned instead of err.Error().
func Handle(c *gin.Context, err error, fallbackMessage string) {
	status := apperrors.HTTPStatus(err)
	message := err.Error()
	if status == http.StatusInternalServerError {
		message = fallbackMessage
	}
	Error(c, status, message)
}
