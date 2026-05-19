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

func TestService_CreateTransaction(t *testing.T) {
	tests := []struct {
		name       string
		req        CreateTransactionRequest
		setup      func(*mockAccountRepository, *mockTransactionRepository)
		wantErrIs  error
		wantAmount float64
	}{
		{
			name: "purchase stores negative amount",
			req:  CreateTransactionRequest{AccountID: 1, OperationTypeID: OperationNormalPurchase, Amount: 100},
			setup: func(ar *mockAccountRepository, tr *mockTransactionRepository) {
				ar.On("GetByID", mock.Anything, int64(1)).Return(&account.Account{ID: 1}, nil)
				tr.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					args[1].(*Transaction).ID = 1
				}).Return(nil)
			},
			wantAmount: -100,
		},
		{
			name: "credit voucher stores positive amount",
			req:  CreateTransactionRequest{AccountID: 1, OperationTypeID: OperationCreditVoucher, Amount: 100},
			setup: func(ar *mockAccountRepository, tr *mockTransactionRepository) {
				ar.On("GetByID", mock.Anything, int64(1)).Return(&account.Account{ID: 1}, nil)
				tr.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					args[1].(*Transaction).ID = 1
				}).Return(nil)
			},
			wantAmount: 100,
		},
		{
			name:      "zero amount rejected",
			req:       CreateTransactionRequest{AccountID: 1, OperationTypeID: OperationNormalPurchase, Amount: 0},
			setup:     func(*mockAccountRepository, *mockTransactionRepository) {},
			wantErrIs: apperrors.ErrValidation,
		},
		{
			name:      "negative amount rejected",
			req:       CreateTransactionRequest{AccountID: 1, OperationTypeID: OperationNormalPurchase, Amount: -50},
			setup:     func(*mockAccountRepository, *mockTransactionRepository) {},
			wantErrIs: apperrors.ErrValidation,
		},
		{
			name:      "invalid operation type",
			req:       CreateTransactionRequest{AccountID: 1, OperationTypeID: 99, Amount: 100},
			setup:     func(*mockAccountRepository, *mockTransactionRepository) {},
			wantErrIs: apperrors.ErrInvalidOperationType,
		},
		{
			name: "account not found",
			req:  CreateTransactionRequest{AccountID: 1, OperationTypeID: OperationNormalPurchase, Amount: 100},
			setup: func(ar *mockAccountRepository, _ *mockTransactionRepository) {
				ar.On("GetByID", mock.Anything, int64(1)).Return(nil, apperrors.ErrNotFound)
			},
			wantErrIs: apperrors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepo := newMockAccountRepository()
			txRepo := newMockTransactionRepository()
			tt.setup(accountRepo, txRepo)

			svc := NewService(accountRepo, txRepo)
			got, err := svc.CreateTransaction(context.Background(), tt.req)

			if tt.wantErrIs != nil {
				require.ErrorIs(t, err, tt.wantErrIs)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantAmount, got.Amount)
			assert.Equal(t, tt.req.OperationTypeID, got.OperationTypeID)
			accountRepo.AssertExpectations(t)
			txRepo.AssertExpectations(t)
		})
	}
}

func TestNormalizeAmount(t *testing.T) {
	tests := []struct {
		name            string
		operationTypeID int
		input           float64
		want            float64
		wantErr         bool
	}{
		{name: "normal purchase", operationTypeID: OperationNormalPurchase, input: 100, want: -100},
		{name: "installment purchase", operationTypeID: OperationInstallmentPurchase, input: 200, want: -200},
		{name: "withdrawal", operationTypeID: OperationWithdrawal, input: 50, want: -50},
		{name: "credit voucher", operationTypeID: OperationCreditVoucher, input: 100, want: 100},
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
