package account

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

	apperrors "github.com/ak2783934/visa-pismo-assessment/pkg/errors"
)

func setupAccountRouter(repo Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)

	svc := NewService(repo)
	handler := NewHandler(svc)

	router := gin.New()
	handler.RegisterRoutes(router.Group("/v1"))
	return router
}

func TestHandler_POST_Accounts_Returns201(t *testing.T) {
	repo := newMockRepository()
	repo.On("GetByDocumentNumber", mock.Anything, "12345678900").Return(nil, apperrors.ErrNotFound)
	repo.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		args[1].(*Account).ID = 1
	}).Return(nil)

	router := setupAccountRouter(repo)

	body, err := json.Marshal(CreateAccountRequest{DocumentNumber: "12345678900"})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/accounts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp Account
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotZero(t, resp.ID)
	assert.Equal(t, "12345678900", resp.DocumentNumber)
	repo.AssertExpectations(t)
}

func TestHandler_GET_Accounts_Returns200(t *testing.T) {
	repo := newMockRepository()
	repo.On("GetByDocumentNumber", mock.Anything, "98765432100").Return(nil, apperrors.ErrNotFound)
	repo.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		args[1].(*Account).ID = 1
	}).Return(nil)
	repo.On("GetByID", mock.Anything, int64(1)).Return(&Account{ID: 1, DocumentNumber: "98765432100"}, nil)

	router := setupAccountRouter(repo)

	createBody, err := json.Marshal(CreateAccountRequest{DocumentNumber: "98765432100"})
	require.NoError(t, err)

	createReq := httptest.NewRequest(http.MethodPost, "/v1/accounts", bytes.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)
	require.Equal(t, http.StatusCreated, createRec.Code)

	var created Account
	require.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &created))

	getReq := httptest.NewRequest(http.MethodGet, "/v1/accounts/"+jsonNumber(created.ID), nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	assert.Equal(t, http.StatusOK, getRec.Code)

	var got Account
	require.NoError(t, json.Unmarshal(getRec.Body.Bytes(), &got))
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, created.DocumentNumber, got.DocumentNumber)
	repo.AssertExpectations(t)
}

func TestHandler_InvalidPayload_Returns400(t *testing.T) {
	repo := newMockRepository()
	router := setupAccountRouter(repo)

	req := httptest.NewRequest(http.MethodPost, "/v1/accounts", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func jsonNumber(n int64) string {
	b, _ := json.Marshal(n)
	return string(b)
}
