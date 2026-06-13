package oca

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mhdiiilham/oca/query"
)

// executor abstracts *sql.DB and *sql.Tx so all repository methods work
// identically inside and outside a transaction.
type executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Repository provides generic CRUD for model T.
// If T implements Tabler its TableName() is used; otherwise the struct name is
// lowercased and an "s" suffix is appended (User → users).
type Repository[T any] struct {
	db      executor
	raw     *sql.DB    // nil when running inside a transaction
	dialect query.Dialect
}

// RepositoryOption configures a Repository at creation time.
type RepositoryOption[T any] func(*Repository[T])

// WithDialect sets the SQL dialect for this repository.
// Defaults to MySQL when not set.
//
//	repo := oca.NewRepository[User](db, oca.WithDialect[User](query.PostgresDialect{}))
func WithDialect[T any](d query.Dialect) RepositoryOption[T] {
	return func(r *Repository[T]) { r.dialect = d }
}

// NewRepository returns a new generic repository backed by db.
func NewRepository[T any](db *sql.DB, opts ...RepositoryOption[T]) GenericStore[T] {
	r := &Repository[T]{db: db, raw: db, dialect: query.MySQLDialect{}}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Tx runs fn inside a single database transaction. The transaction is
// committed when fn returns nil, or rolled back on any error.
// Calling Tx on a repository that is already inside a transaction returns an error.
func (r *Repository[T]) Tx(ctx context.Context, fn func(store GenericStore[T]) error) error {
	if r.raw == nil {
		return fmt.Errorf("oca: Tx cannot be nested inside an existing transaction")
	}
	tx, err := r.raw.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("oca: begin transaction: %w", err)
	}
	txRepo := &Repository[T]{db: tx, raw: nil, dialect: r.dialect}
	if err := fn(txRepo); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("oca: commit transaction: %w", err)
	}
	return nil
}
