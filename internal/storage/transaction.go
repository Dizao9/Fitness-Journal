package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type txKey struct{}

type Manager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *Manager {
	return &Manager{
		db: db,
	}
}

func (m *Manager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	ctxWithTX := context.WithValue(ctx, txKey{}, tx)

	err = fn(ctxWithTX)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback transaction: %w (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

type dbOps interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func getQueryEngine(ctx context.Context, db *sql.DB) dbOps {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	if ok && tx != nil {
		return tx
	}
	return db
}
