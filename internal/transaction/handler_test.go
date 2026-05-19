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
	apperrors "github.com/ak2783934/visa-pismo-assessment/pkg/errors"
)

func setupTransactionRouter(accountRepo account.Repository, txRepo Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc := NewService(accountRepo, txRepo)
	handler := NewHandler(svc)
	router := gin.New()
	handler.RegisterRoutes(router.Group("/v1"))
	return router
}

func mustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func TestHandler_CreateTransaction(t *testing.T) {
	tests := []struct {
		name       string
		body       []byte
		setup      func(*mockAccountRepository, *mockTransactionRepository)
		wantStatus int
		checkBody  func(*testing.T, []byte)
	}{
		{
			name: "success",
			body: mustMarshal(CreateTransactionRequest{AccountID: 1, OperationTypeID: OperationNormalPurchase, Amount: 100}),
			setup: func(ar *mockAccountRepository, tr *mockTransactionRepository) {
				ar.On("GetByID", mock.Anything, int64(1)).Return(&account.Account{ID: 1}, nil)
				tr.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					args[1].(*Transaction).ID = 1
				}).Return(nil)
			},
			wantStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body []byte) {
				var resp Transaction
				require.NoError(t, json.Unmarshal(body, &resp))
				assert.NotZero(t, resp.ID)
				assert.Equal(t, int64(1), resp.AccountID)
				assert.Equal(t, -100.0, resp.Amount)
			},
		},
		{
			name:       "invalid JSON",
			body:       []byte(`{`),
			setup:      func(*mockAccountRepository, *mockTransactionRepository) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "account not found",
			body: mustMarshal(CreateTransactionRequest{AccountID: 1, OperationTypeID: OperationNormalPurchase, Amount: 100}),
			setup: func(ar *mockAccountRepository, _ *mockTransactionRepository) {
				ar.On("GetByID", mock.Anything, int64(1)).Return(nil, apperrors.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid operation type",
			body:       mustMarshal(CreateTransactionRequest{AccountID: 1, OperationTypeID: 99, Amount: 100}),
			setup:      func(*mockAccountRepository, *mockTransactionRepository) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "zero amount",
			body:       mustMarshal(CreateTransactionRequest{AccountID: 1, OperationTypeID: OperationNormalPurchase, Amount: 0}),
			setup:      func(*mockAccountRepository, *mockTransactionRepository) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepo := newMockAccountRepository()
			txRepo := newMockTransactionRepository()
			tt.setup(accountRepo, txRepo)

			router := setupTransactionRouter(accountRepo, txRepo)
			req := httptest.NewRequest(http.MethodPost, "/v1/transactions", bytes.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.checkBody != nil {
				tt.checkBody(t, w.Body.Bytes())
			}
			accountRepo.AssertExpectations(t)
			txRepo.AssertExpectations(t)
		})
	}
}
