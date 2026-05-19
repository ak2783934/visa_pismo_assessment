package transaction

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ak2783934/visa-pismo-assessment/internal/account"
)

func setupTransactionRouter(accountRepo account.Repository, txRepo Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)

	svc := NewService(accountRepo, txRepo)
	handler := NewHandler(svc)

	router := gin.New()
	handler.RegisterRoutes(router.Group("/v1"))
	return router
}

func TestHandler_POST_Transactions_Returns201(t *testing.T) {
	accountRepo := newMockAccountRepository()
	accountRepo.On("GetByID", mock.Anything, int64(1)).
		Return(&account.Account{ID: 1, DocumentNumber: "test"}, nil)

	txRepo := newMockTransactionRepository()
	txRepo.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		args[1].(*Transaction).ID = 1
	}).Return(nil)

	router := setupTransactionRouter(accountRepo, txRepo)

	body, err := json.Marshal(CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: OperationNormalPurchase,
		Amount:          100,
	})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp Transaction
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotZero(t, resp.ID)
	assert.Equal(t, int64(1), resp.AccountID)
	assert.Equal(t, -100.0, resp.Amount)

	accountRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestHandler_InvalidPayload_Returns400(t *testing.T) {
	accountRepo := newMockAccountRepository()
	txRepo := newMockTransactionRepository()

	router := setupTransactionRouter(accountRepo, txRepo)

	req := httptest.NewRequest(http.MethodPost, "/v1/transactions", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
