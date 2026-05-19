package transaction

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, tx *Transaction) error
}

type sqliteRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqliteRepository{db: db}
}

func (r *sqliteRepository) Create(ctx context.Context, tx *Transaction) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO transactions (account_id, operation_type_id, amount, event_date)
		 VALUES (?, ?, ?, ?)`,
		tx.AccountID,
		tx.OperationTypeID,
		tx.Amount,
		tx.EventDate,
	)
	if err != nil {
		return fmt.Errorf("insert transaction: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}

	tx.ID = id
	return nil
}
