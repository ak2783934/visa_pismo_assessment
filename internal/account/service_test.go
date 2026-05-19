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

func TestService_CreateAccount_Success(t *testing.T) {
	repo := newMockRepository()
	repo.On("GetByDocumentNumber", mock.Anything, "12345678900").Return(nil, apperrors.ErrNotFound)
	repo.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		args[1].(*Account).ID = 1
	}).Return(nil)

	svc := NewService(repo)
	account, err := svc.CreateAccount(context.Background(), CreateAccountRequest{DocumentNumber: "12345678900"})
	require.NoError(t, err)
	assert.NotZero(t, account.ID)
	assert.Equal(t, "12345678900", account.DocumentNumber)
	repo.AssertExpectations(t)
}

func TestService_CreateAccount_EmptyDocumentNumber(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()

	_, err := svc.CreateAccount(ctx, CreateAccountRequest{DocumentNumber: ""})
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrValidation)

	_, err = svc.CreateAccount(ctx, CreateAccountRequest{DocumentNumber: "   "})
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrValidation)
}

func TestService_CreateAccount_DuplicateDocumentNumber(t *testing.T) {
	repo := newMockRepository()
	repo.On("GetByDocumentNumber", mock.Anything, "12345678900").
		Return(&Account{ID: 1, DocumentNumber: "12345678900"}, nil)

	svc := NewService(repo)
	_, err := svc.CreateAccount(context.Background(), CreateAccountRequest{DocumentNumber: "12345678900"})
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrDuplicate)
	repo.AssertExpectations(t)
}

func TestService_GetAccount_Success(t *testing.T) {
	repo := newMockRepository()
	repo.On("GetByDocumentNumber", mock.Anything, "12345678900").Return(nil, apperrors.ErrNotFound)
	repo.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		args[1].(*Account).ID = 1
	}).Return(nil)
	repo.On("GetByID", mock.Anything, int64(1)).Return(&Account{ID: 1, DocumentNumber: "12345678900"}, nil)

	svc := NewService(repo)
	ctx := context.Background()

	created, err := svc.CreateAccount(ctx, CreateAccountRequest{DocumentNumber: "12345678900"})
	require.NoError(t, err)

	account, err := svc.GetAccount(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, account.ID)
	assert.Equal(t, created.DocumentNumber, account.DocumentNumber)
	repo.AssertExpectations(t)
}

func TestService_GetAccount_NotFound(t *testing.T) {
	repo := newMockRepository()
	repo.On("GetByID", mock.Anything, int64(999)).Return(nil, apperrors.ErrNotFound)

	svc := NewService(repo)
	_, err := svc.GetAccount(context.Background(), 999)
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
	repo.AssertExpectations(t)
}
