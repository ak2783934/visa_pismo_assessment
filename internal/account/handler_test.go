package account

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func mustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func TestHandler_CreateAccount(t *testing.T) {
	tests := []struct {
		name       string
		body       []byte
		setup      func(*mockRepository)
		wantStatus int
		checkBody  func(*testing.T, []byte)
	}{
		{
			name: "success",
			body: mustMarshal(CreateAccountRequest{DocumentNumber: "12345678900"}),
			setup: func(r *mockRepository) {
				r.On("GetByDocumentNumber", mock.Anything, "12345678900").Return(nil, apperrors.ErrNotFound)
				r.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					args[1].(*Account).ID = 1
				}).Return(nil)
			},
			wantStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body []byte) {
				var resp Account
				require.NoError(t, json.Unmarshal(body, &resp))
				assert.NotZero(t, resp.ID)
				assert.Equal(t, "12345678900", resp.DocumentNumber)
			},
		},
		{
			name:       "invalid JSON",
			body:       []byte(`{`),
			setup:      func(*mockRepository) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate document number",
			body: mustMarshal(CreateAccountRequest{DocumentNumber: "12345678900"}),
			setup: func(r *mockRepository) {
				r.On("GetByDocumentNumber", mock.Anything, "12345678900").
					Return(&Account{ID: 1, DocumentNumber: "12345678900"}, nil)
			},
			wantStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepository()
			tt.setup(repo)

			router := setupAccountRouter(repo)
			req := httptest.NewRequest(http.MethodPost, "/v1/accounts", bytes.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.checkBody != nil {
				tt.checkBody(t, w.Body.Bytes())
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestHandler_GetAccount(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setup      func(*mockRepository)
		wantStatus int
		checkBody  func(*testing.T, []byte)
	}{
		{
			name: "success",
			id:   "1",
			setup: func(r *mockRepository) {
				r.On("GetByID", mock.Anything, int64(1)).Return(&Account{ID: 1, DocumentNumber: "12345678900"}, nil)
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp Account
				require.NoError(t, json.Unmarshal(body, &resp))
				assert.Equal(t, int64(1), resp.ID)
				assert.Equal(t, "12345678900", resp.DocumentNumber)
			},
		},
		{
			name: "not found",
			id:   "999",
			setup: func(r *mockRepository) {
				r.On("GetByID", mock.Anything, int64(999)).Return(nil, apperrors.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid id",
			id:         "abc",
			setup:      func(*mockRepository) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepository()
			tt.setup(repo)

			router := setupAccountRouter(repo)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/accounts/%s", tt.id), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.checkBody != nil {
				tt.checkBody(t, w.Body.Bytes())
			}
			repo.AssertExpectations(t)
		})
	}
}
