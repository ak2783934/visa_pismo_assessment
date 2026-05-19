package transaction

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ak2783934/visa-pismo-assessment/internal/account"
	apperrors "github.com/ak2783934/visa-pismo-assessment/pkg/errors"
)

type mockAccountRepository struct {
	*mock.Mock
}

func (m *mockAccountRepository) Create(ctx context.Context, a *account.Account) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *mockAccountRepository) GetByID(ctx context.Context, id int64) (*account.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*account.Account), args.Error(1)
}

func (m *mockAccountRepository) GetByDocumentNumber(ctx context.Context, documentNumber string) (*account.Account, error) {
	args := m.Called(ctx, documentNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*account.Account), args.Error(1)
}

func newMockAccountRepository() *mockAccountRepository {
	return &mockAccountRepository{Mock: &mock.Mock{}}
}

type mockTransactionRepository struct {
	*mock.Mock
}

func (m *mockTransactionRepository) Create(ctx context.Context, tx *Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func newMockTransactionRepository() *mockTransactionRepository {
	return &mockTransactionRepository{Mock: &mock.Mock{}}
}

func TestService_CreateTransaction_Operation1NegativeAmount(t *testing.T) {
	accountRepo := newMockAccountRepository()
	accountRepo.On("GetByID", mock.Anything, int64(1)).
		Return(&account.Account{ID: 1, DocumentNumber: "test"}, nil)

	txRepo := newMockTransactionRepository()
	txRepo.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		args[1].(*Transaction).ID = 1
	}).Return(nil)

	svc := NewService(accountRepo, txRepo)
	tx, err := svc.CreateTransaction(context.Background(), CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: OperationNormalPurchase,
		Amount:          100,
	})
	require.NoError(t, err)
	assert.Equal(t, -100.0, tx.Amount)
	assert.Equal(t, OperationNormalPurchase, tx.OperationTypeID)
	accountRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestService_CreateTransaction_Operation4PositiveAmount(t *testing.T) {
	accountRepo := newMockAccountRepository()
	accountRepo.On("GetByID", mock.Anything, int64(1)).
		Return(&account.Account{ID: 1, DocumentNumber: "test"}, nil)

	txRepo := newMockTransactionRepository()
	txRepo.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		args[1].(*Transaction).ID = 1
	}).Return(nil)

	svc := NewService(accountRepo, txRepo)
	tx, err := svc.CreateTransaction(context.Background(), CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: OperationCreditVoucher,
		Amount:          100,
	})
	require.NoError(t, err)
	assert.Equal(t, 100.0, tx.Amount)
	assert.Equal(t, OperationCreditVoucher, tx.OperationTypeID)
	accountRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestService_CreateTransaction_InvalidOperationType(t *testing.T) {
	accountRepo := newMockAccountRepository()
	txRepo := newMockTransactionRepository()

	svc := NewService(accountRepo, txRepo)
	_, err := svc.CreateTransaction(context.Background(), CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 99,
		Amount:          100,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrInvalidOperationType)
}

func TestService_CreateTransaction_MissingAccount(t *testing.T) {
	accountRepo := newMockAccountRepository()
	accountRepo.On("GetByID", mock.Anything, int64(1)).
		Return(nil, apperrors.ErrNotFound)

	txRepo := newMockTransactionRepository()

	svc := NewService(accountRepo, txRepo)
	_, err := svc.CreateTransaction(context.Background(), CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: OperationNormalPurchase,
		Amount:          100,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
	accountRepo.AssertExpectations(t)
}

func TestService_CreateTransaction_AmountNotPositive(t *testing.T) {
	accountRepo := newMockAccountRepository()
	txRepo := newMockTransactionRepository()

	svc := NewService(accountRepo, txRepo)
	ctx := context.Background()

	_, err := svc.CreateTransaction(ctx, CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: OperationNormalPurchase,
		Amount:          0,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrValidation)

	_, err = svc.CreateTransaction(ctx, CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: OperationNormalPurchase,
		Amount:          -50,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrValidation)
}

func TestNormalizeAmount(t *testing.T) {
	tests := []struct {
		name            string
		operationTypeID int
		input           float64
		want            float64
		wantErr         bool
	}{
		{name: "normal purchase positive", operationTypeID: 1, input: 100, want: -100},
		{name: "credit voucher positive", operationTypeID: 4, input: 100, want: 100},
		{name: "invalid operation type", operationTypeID: 99, input: 100, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeAmount(tt.operationTypeID, tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
