package oca

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mhdiiilham/oca/query"
)

// FindFilter carries the options applied to a SELECT query.
type FindFilter struct {
	Where  []query.Condition
	Order  string
	Limit  int
	Offset int
}

// FilterOptions modifies a FindFilter.
type FilterOptions func(*FindFilter)

// Where adds a WHERE condition.
//
//	Where(query.C("age").Gt(18))
func Where(cond query.Condition) FilterOptions {
	return func(ff *FindFilter) {
		ff.Where = append(ff.Where, cond)
	}
}

// Direction is the sort direction for OrderBy.
type Direction string

const (
	ASC  Direction = "ASC"
	DESC Direction = "DESC"
)

// OrderBy sets a safe ORDER BY on a column and direction.
//
//	oca.OrderBy("created_at", oca.DESC)
func OrderBy(col string, dir Direction) FilterOptions {
	return func(ff *FindFilter) {
		ff.Order = col + " " + string(dir)
	}
}

// OrderByRaw sets a raw ORDER BY clause. The string is written directly into
// SQL without escaping — do not pass user input here.
func OrderByRaw(raw string) FilterOptions {
	return func(ff *FindFilter) { ff.Order = raw }
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

// Finds retrieves records matching the provided filter options.
func (r *Repository[T]) Finds(ctx context.Context, opts ...FilterOptions) ([]T, error) {
	var entity T
	builder := query.From(resolveTableName(entity)).WithDialect(r.dialect).Select(getColumnNames(entity)...)

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
		return nil, fmt.Errorf("oca: finds: %w", err)
	}
	defer rows.Close()

	return scanRows[T](rows)
}

// FindOne retrieves one record matching the filter. Returns sql.ErrNoRows when nothing matches.
func (r *Repository[T]) FindOne(ctx context.Context, opts ...FilterOptions) (*T, error) {
	opts = append(opts, Limit(1))
	results, err := r.Finds(ctx, opts...)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}
	return &results[0], nil
}

// applyFilters applies WHERE, ORDER BY, LIMIT, OFFSET to the query builder.
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
