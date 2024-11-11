package repository

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Intiqo/app-platform/internal/domain"
)

type pgxSettingRepository struct {
	db  *pgxpool.Pool
	sqt sq.StatementBuilderType
}

// NewSettingRepository creates a new setting repository
func NewSettingRepository(db *pgxpool.Pool) domain.SettingRepository {
	return &pgxSettingRepository{
		db:  db,
		sqt: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *pgxSettingRepository) FindByID(ctx context.Context, id uuid.UUID) (result domain.Setting, err error) {
	// Check if the context has a transaction
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Construct the query
	q := `SELECT * FROM settings WHERE id = $1 AND deleted_at IS NULL`
	args := []interface{}{id}

	// Execute the query
	var rows pgx.Rows
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		rows, err = tx.Query(ctx, q, args...)
	} else {
		rows, err = r.db.Query(ctx, q, args...)
	}
	defer rows.Close()
	if err != nil {
		return result, err
	}
	if rows == nil {
		return result, domain.DataNotFoundError{}
	}

	// Collect the results
	result, err = pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[domain.Setting])
	if err != nil {
		return result, err
	}

	// Return the result
	return result, nil
}

func (r *pgxSettingRepository) Filter(ctx context.Context, in domain.FilterSettingsByCriteriaInput, opts domain.QueryOptions) (result []domain.Setting, total int64, err error) {
	// Check if the context has a transaction
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Build the criteria
	f := r.sqt

	if len(in.Keys) > 0 {
		f = f.Where(sq.Eq{"key": in.Keys})
	}

	f = f.Where("deleted_at IS NULL")

	// Build the count query
	cb := f.Select("COUNT(*)").
		From("settings")
	cq, cargs, err := cb.ToSql()
	if err != nil {
		return result, total, err
	}

	// Execute the count query
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		err = tx.QueryRow(ctx, cq, cargs...).Scan(&total)
	} else {
		err = r.db.QueryRow(ctx, cq, cargs...).Scan(&total)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return result, total, nil
		}
		return result, total, err
	}

	// Build the query
	qb := f.Select("*").
		From("settings")

	if opts.Limit > 0 {
		qb = qb.Limit(uint64(opts.Limit))
	}
	if opts.Offset > 0 {
		qb = qb.Offset(uint64(opts.Offset))
	}
	dq, dargs, err := qb.ToSql()
	if err != nil {
		return result, total, err
	}

	var rows pgx.Rows
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		rows, err = tx.Query(ctx, dq, dargs...)
	} else {
		rows, err = r.db.Query(ctx, dq, dargs...)
	}
	defer rows.Close()
	if err != nil {
		return result, total, err
	}
	if rows == nil {
		return result, total, nil
	}

	// Collect the data into the result
	result, err = pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.Setting])
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return result, total, nil
	}

	return result, total, err
}

func (r *pgxSettingRepository) Create(ctx context.Context, entity *domain.Setting) (err error) {
	// Check if the context has a transaction
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Construct the query
	q := `INSERT INTO settings (key, value) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	args := []interface{}{entity.Key, entity.Value}

	// Execute the query
	var row pgx.Row
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		row = tx.QueryRow(ctx, q, args...)
	} else {
		row = r.db.QueryRow(ctx, q, args...)
	}

	// Collect the result
	err = row.Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)
	if err != nil {
		return err
	}

	// Return the result
	return err
}

func (r *pgxSettingRepository) CreateMultiple(ctx context.Context, entities []*domain.Setting) (err error) {
	// Check if the context has a transaction
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Create a batch
	b := &pgx.Batch{}

	// Add queries to the batch
	q := `INSERT INTO settings (key, value) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	for idx, entity := range entities {
		// Create the data
		args := []interface{}{entity.Key, entity.Value}
		b.Queue(q, args...).QueryRow(func(row pgx.Row) error {
			return row.Scan(&entities[idx].ID, &entities[idx].CreatedAt, &entities[idx].UpdatedAt)
		})
	}

	// Execute the batch
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		err = tx.SendBatch(ctx, b).Close()
	} else {
		err = r.db.SendBatch(ctx, b).Close()
	}

	// Return the result
	return err
}

func (r *pgxSettingRepository) Update(ctx context.Context, entity *domain.Setting) (err error) {
	// Check if the context has a transaction
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Construct the query
	q := `UPDATE settings SET key = $1, value = $2, updated_at = NOW() WHERE id = $3 RETURNING updated_at`
	args := []interface{}{entity.Key, entity.Value, entity.ID}

	// Execute the query
	var row pgx.Row
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		row = tx.QueryRow(ctx, q, args...)
	} else {
		row = r.db.QueryRow(ctx, q, args...)
	}

	// Collect the result
	err = row.Scan(&entity.UpdatedAt)
	if err != nil {
		return err
	}

	// Return the result
	return err
}

func (r *pgxSettingRepository) UpdateMultiple(ctx context.Context, entities []*domain.Setting) (err error) {
	// Check if the context has a transaction
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Create a batch
	b := &pgx.Batch{}

	// Add queries to the batch
	q := `UPDATE settings SET key = $1, value = $2, updated_at = NOW() WHERE id = $3 RETURNING updated_at`
	for idx, entity := range entities {
		// Create the data
		args := []interface{}{entity.Key, entity.Value, entity.ID}
		b.Queue(q, args...).QueryRow(func(row pgx.Row) error {
			return row.Scan(&entities[idx].UpdatedAt)
		})
	}

	// Execute the batch
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		err = tx.SendBatch(ctx, b).Close()
	} else {
		err = r.db.SendBatch(ctx, b).Close()
	}

	// Return the result
	return err
}

func (r *pgxSettingRepository) DeleteByID(ctx context.Context, id uuid.UUID) (err error) {
	// Check if the context has a transaction
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Construct the query
	q := `UPDATE settings SET deleted_at = NOW() WHERE id = $1`
	args := []interface{}{id}

	// Execute the query
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		_, err = tx.Exec(ctx, q, args...)
	} else {
		_, err = r.db.Exec(ctx, q, args...)
	}

	// Return the result
	return err
}

func (r *pgxSettingRepository) DeleteByIDs(ctx context.Context, ids []uuid.UUID) (err error) {
	// Check if the context has a transaction
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Construct the query
	q := `UPDATE settings SET deleted_at = NOW() WHERE id = ANY($1)`
	args := []interface{}{ids}

	// Execute the query
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		_, err = tx.Exec(ctx, q, args...)
	} else {
		_, err = r.db.Exec(ctx, q, args...)
	}

	// Return the result
	return err
}
