package account

import (
	"net/http"
	"strconv"

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
	r.POST("/accounts", h.CreateAccount)
	r.GET("/accounts/:id", h.GetAccount)
}

// CreateAccount creates a new account.
//
//	@Summary		Create account
//	@Description	Create an account with a unique document number
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateAccountRequest	true	"Create account request"
//	@Success		201		{object}	Account
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		409		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/accounts [post]
func (h *Handler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Handle(c, apperrors.Wrap(apperrors.ErrValidation, err.Error()), "")
		return
	}

	account, err := h.service.CreateAccount(c.Request.Context(), req)
	if err != nil {
		response.Handle(c, err, "failed to create account")
		return
	}

	response.JSON(c, http.StatusCreated, account)
}

// GetAccount returns an account by ID.
//
//	@Summary		Get account
//	@Description	Get account details by account ID
//	@Tags			accounts
//	@Produce		json
//	@Param			id	path		int	true	"Account ID"
//	@Success		200	{object}	Account
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/accounts/{id} [get]
func (h *Handler) GetAccount(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Handle(c, apperrors.Wrap(apperrors.ErrValidation, "invalid account id"), "")
		return
	}

	account, err := h.service.GetAccount(c.Request.Context(), id)
	if err != nil {
		response.Handle(c, err, "failed to get account")
		return
	}

	response.JSON(c, http.StatusOK, account)
}
