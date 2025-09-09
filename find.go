package oca

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mhdiiilham/oca/query"
)

// FindFilter defines filters for the Find method.
type FindFilter struct {
	Where  []query.Condition // WHERE conditions (multiple)
	Order  string            // ORDER BY clause
	Limit  int               // LIMIT
	Offset int               // OFFSET
}

// FilterOptions modifies a FindFilter.
type FilterOptions func(*FindFilter)

// Where adds a WHERE condition using the query DSL.
//
// Example:
//
//	Where(query.C("age").Gt(18))
func Where(cond query.Condition) FilterOptions {
	return func(ff *FindFilter) {
		ff.Where = append(ff.Where, cond)
	}
}

// OrderBy sets the ORDER BY clause.
func OrderBy(order string) FilterOptions {
	return func(ff *FindFilter) { ff.Order = order }
}

// Limit sets the LIMIT clause.
func Limit(limit int) FilterOptions {
	return func(ff *FindFilter) { ff.Limit = limit }
}

// Offset sets the OFFSET clause.
func Offset(offset int) FilterOptions {
	return func(ff *FindFilter) { ff.Offset = offset }
}

func defaultFilter() FindFilter {
	return FindFilter{}
}

// Find retrieves multiple records from the database matching the filter.
// It automatically maps database rows to struct fields based on the `db` tags.
// Find retrieves records from the database based on provided filter options.
// It returns a slice of T or an error.
func (r *Repository[T]) Find(ctx context.Context, opts ...FilterOptions) ([]T, error) {
	var entity T
	// Build the SQL query
	builder := query.From(resolveTableName(entity)).Select(getColumnNames(entity)...)

	filter := defaultFilter()
	for _, opt := range opts {
		if opt != nil {
			opt(&filter)
		}
	}

	applyFilters(builder, filter)

	sqlStr, args := builder.Build()
	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	return scanRows[T](rows)
}

// FindOne retrieves a single record from the database based on the filter.
// Returns sql.ErrNoRows if nothing is found.
func (r *Repository[T]) FindOne(ctx context.Context, opts ...FilterOptions) (*T, error) {
	opts = append(opts, Limit(1))
	results, err := r.Find(ctx, opts...)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}

	return &results[0], nil
}

// applyFilters applies WHERE, ORDER BY, LIMIT, OFFSET to the query builder
func applyFilters(b *query.Builder, f FindFilter) {
	for _, cond := range f.Where {
		b.Where(cond)
	}
	if f.Order != "" {
		b.OrderBy(f.Order)
	}
	if f.Limit > 0 {
		b.Limit(f.Limit)
	}
	if f.Offset > 0 {
		b.Offset(f.Offset)
	}
}
