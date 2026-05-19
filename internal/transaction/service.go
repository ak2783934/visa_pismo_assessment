package transaction

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/ak2783934/visa-pismo-assessment/internal/account"
	apperrors "github.com/ak2783934/visa-pismo-assessment/pkg/errors"
)

const (
	OperationNormalPurchase      = 1
	OperationInstallmentPurchase = 2
	OperationWithdrawal          = 3
	OperationCreditVoucher       = 4
)

type Service struct {
	accountRepo account.Repository
	txRepo      Repository
}

func NewService(accountRepo account.Repository, txRepo Repository) *Service {
	return &Service{
		accountRepo: accountRepo,
		txRepo:      txRepo,
	}
}

func (s *Service) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*Transaction, error) {
	if req.Amount <= 0 {
		return nil, apperrors.Wrap(apperrors.ErrValidation, "amount must be greater than 0")
	}

	if !isValidOperationType(req.OperationTypeID) {
		return nil, apperrors.Wrap(apperrors.ErrInvalidOperationType, "operation type must be between 1 and 4")
	}

	if _, err := s.accountRepo.GetByID(ctx, req.AccountID); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.Wrap(apperrors.ErrNotFound, "account not found")
		}
		return nil, err
	}

	amount, err := normalizeAmount(req.OperationTypeID, req.Amount)
	if err != nil {
		return nil, err
	}

	tx := &Transaction{
		AccountID:       req.AccountID,
		OperationTypeID: req.OperationTypeID,
		Amount:          amount,
		EventDate:       time.Now().UTC(),
	}

	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func isValidOperationType(operationTypeID int) bool {
	return operationTypeID >= OperationNormalPurchase && operationTypeID <= OperationCreditVoucher
}

func normalizeAmount(operationTypeID int, amount float64) (float64, error) {
	magnitude := math.Abs(amount)

	switch operationTypeID {
	case OperationNormalPurchase, OperationInstallmentPurchase, OperationWithdrawal:
		return -magnitude, nil
	case OperationCreditVoucher:
		return magnitude, nil
	default:
		return 0, apperrors.Wrap(apperrors.ErrInvalidOperationType, "invalid operation type")
	}
}
