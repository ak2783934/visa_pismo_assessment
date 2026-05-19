package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	apperrors "github.com/ak2783934/visa-pismo-assessment/pkg/errors"
)

func isUniqueConstraint(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "unique constraint")
}

type Repository interface {
	Create(ctx context.Context, account *Account) error
	GetByID(ctx context.Context, id int64) (*Account, error)
	GetByDocumentNumber(ctx context.Context, documentNumber string) (*Account, error)
}

type sqliteRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqliteRepository{db: db}
}

func (r *sqliteRepository) Create(ctx context.Context, account *Account) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO accounts (document_number) VALUES (?)`,
		account.DocumentNumber,
	)
	if err != nil {
		if isUniqueConstraint(err) {
			return apperrors.Wrap(apperrors.ErrDuplicate, "account with this document number already exists")
		}
		return fmt.Errorf("insert account: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}

	account.ID = id
	return nil
}

func (r *sqliteRepository) GetByID(ctx context.Context, id int64) (*Account, error) {
	return r.getOne(ctx,
		`SELECT id, document_number FROM accounts WHERE id = ?`,
		id,
	)
}

func (r *sqliteRepository) GetByDocumentNumber(ctx context.Context, documentNumber string) (*Account, error) {
	return r.getOne(ctx,
		`SELECT id, document_number FROM accounts WHERE document_number = ?`,
		documentNumber,
	)
}

func (r *sqliteRepository) getOne(ctx context.Context, query string, arg any) (*Account, error) {
	row := r.db.QueryRowContext(ctx, query, arg)

	var account Account
	if err := row.Scan(&account.ID, &account.DocumentNumber); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.Wrap(apperrors.ErrNotFound, "account not found")
		}
		return nil, fmt.Errorf("scan account: %w", err)
	}

	return &account, nil
}
