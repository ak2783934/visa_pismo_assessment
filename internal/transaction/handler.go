package transaction

import (
	"net/http"

	apperrors "github.com/ak2783934/visa-pismo-assessment/pkg/errors"
	"github.com/ak2783934/visa-pismo-assessment/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r gin.IRoutes) {
	r.POST("/transactions", h.CreateTransaction)
}

// CreateTransaction records a new transaction.
//
//	@Summary		Create transaction
//	@Description	Record a transaction. Amount sign is normalized by operation type (1-3 negative, 4 positive).
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateTransactionRequest	true	"Create transaction request"
//	@Success		201		{object}	Transaction
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		404		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/transactions [post]
func (h *Handler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Handle(c, apperrors.Wrap(apperrors.ErrValidation, err.Error()), "")
		return
	}

	tx, err := h.service.CreateTransaction(c.Request.Context(), req)
	if err != nil {
		response.Handle(c, err, "failed to create transaction")
		return
	}

	response.JSON(c, http.StatusCreated, tx)
}
