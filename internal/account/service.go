package account

import (
	"context"
	"errors"
	"strings"

	apperrors "github.com/ak2783934/visa-pismo-assessment/pkg/errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAccount(ctx context.Context, req CreateAccountRequest) (*Account, error) {
	documentNumber := strings.TrimSpace(req.DocumentNumber)
	if documentNumber == "" {
		return nil, apperrors.Wrap(apperrors.ErrValidation, "document number is required")
	}

	existing, err := s.repo.GetByDocumentNumber(ctx, documentNumber)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, apperrors.Wrap(apperrors.ErrDuplicate, "account with this document number already exists")
	}

	account := &Account{DocumentNumber: documentNumber}
	if err := s.repo.Create(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *Service) GetAccount(ctx context.Context, id int64) (*Account, error) {
	return s.repo.GetByID(ctx, id)
}
