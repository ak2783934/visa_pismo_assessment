package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	apperrors "github.com/ak2783934/visa-pismo-assessment/pkg/errors"
)

type mockRepository struct {
	*mock.Mock
}

func (m *mockRepository) Create(ctx context.Context, a *Account) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *mockRepository) GetByID(ctx context.Context, id int64) (*Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Account), args.Error(1)
}

func (m *mockRepository) GetByDocumentNumber(ctx context.Context, documentNumber string) (*Account, error) {
	args := m.Called(ctx, documentNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Account), args.Error(1)
}

func newMockRepository() *mockRepository {
	return &mockRepository{Mock: &mock.Mock{}}
}

func TestService_CreateAccount(t *testing.T) {
	tests := []struct {
		name      string
		req       CreateAccountRequest
		setup     func(*mockRepository)
		wantErrIs error
	}{
		{
			name: "success",
			req:  CreateAccountRequest{DocumentNumber: "12345678900"},
			setup: func(r *mockRepository) {
				r.On("GetByDocumentNumber", mock.Anything, "12345678900").Return(nil, apperrors.ErrNotFound)
				r.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					args[1].(*Account).ID = 1
				}).Return(nil)
			},
		},
		{
			name:      "empty document number",
			req:       CreateAccountRequest{DocumentNumber: ""},
			setup:     func(*mockRepository) {},
			wantErrIs: apperrors.ErrValidation,
		},
		{
			name:      "whitespace document number",
			req:       CreateAccountRequest{DocumentNumber: "   "},
			setup:     func(*mockRepository) {},
			wantErrIs: apperrors.ErrValidation,
		},
		{
			name: "duplicate document number",
			req:  CreateAccountRequest{DocumentNumber: "12345678900"},
			setup: func(r *mockRepository) {
				r.On("GetByDocumentNumber", mock.Anything, "12345678900").
					Return(&Account{ID: 1, DocumentNumber: "12345678900"}, nil)
			},
			wantErrIs: apperrors.ErrDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepository()
			tt.setup(repo)

			svc := NewService(repo)
			got, err := svc.CreateAccount(context.Background(), tt.req)

			if tt.wantErrIs != nil {
				require.ErrorIs(t, err, tt.wantErrIs)
				return
			}
			require.NoError(t, err)
			assert.NotZero(t, got.ID)
			assert.Equal(t, tt.req.DocumentNumber, got.DocumentNumber)
			repo.AssertExpectations(t)
		})
	}
}

func TestService_GetAccount(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		setup     func(*mockRepository)
		wantErrIs error
	}{
		{
			name: "success",
			id:   1,
			setup: func(r *mockRepository) {
				r.On("GetByID", mock.Anything, int64(1)).Return(&Account{ID: 1, DocumentNumber: "12345678900"}, nil)
			},
		},
		{
			name: "not found",
			id:   999,
			setup: func(r *mockRepository) {
				r.On("GetByID", mock.Anything, int64(999)).Return(nil, apperrors.ErrNotFound)
			},
			wantErrIs: apperrors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepository()
			tt.setup(repo)

			svc := NewService(repo)
			got, err := svc.GetAccount(context.Background(), tt.id)

			if tt.wantErrIs != nil {
				require.ErrorIs(t, err, tt.wantErrIs)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.id, got.ID)
			repo.AssertExpectations(t)
		})
	}
}
