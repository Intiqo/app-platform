package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Intiqo/app-platform/internal/domain"
)

type TxKeyType string

const TxKey TxKeyType = "App-Transactioner"

type transactioner struct {
	db *pgxpool.Pool
}

func NewTransactioner(db *pgxpool.Pool) domain.Transactioner {
	return &transactioner{
		db: db,
	}
}

func (t *transactioner) Begin(ctx context.Context) (result context.Context, err error) {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return result, err
	}
	result = context.WithValue(ctx, TxKey, tx)
	return result, nil
}

func (t *transactioner) Commit(ctx context.Context) (err error) {
	tx, ok := ctx.Value(TxKey).(*pgxpool.Tx)
	if !ok {
		return ErrTransactionNotFound
	}
	return tx.Commit(ctx)
}

func (t *transactioner) Rollback(ctx context.Context, err error) {
	if err == nil {
		return
	}
	tx, ok := ctx.Value(TxKey).(*pgxpool.Tx)
	if !ok {
		slog.Error("no transaction found")
	}
	err = tx.Rollback(ctx)
	if err != nil {
		slog.Error("failed to rollback transaction", "error", err)
	}
}

var ErrTransactionNotFound = errors.New("no transaction found")
