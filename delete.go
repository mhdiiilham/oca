package oca

import (
	"context"
	"fmt"

	"github.com/mhdiiilham/oca/query"
)

// Delete removes rows matching the filter. At least one Where() condition is
// required; call DeleteAll to intentionally remove every row.
func (r *Repository[T]) Delete(ctx context.Context, opts ...FilterOptions) error {
	filter := defaultFilter()
	for _, opt := range opts {
		if opt != nil {
			opt(&filter)
		}
	}

	if len(filter.Where) == 0 {
		return fmt.Errorf("oca: Delete requires at least one Where condition to prevent accidental full-table deletion; use DeleteAll to delete every row")
	}

	var t T
	builder := query.Delete(resolveTableName(t)).WithDialect(r.dialect)
	builder.Where(filter.Where...)

	sqlStr, args := builder.Build()
	if _, err := r.db.ExecContext(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("oca: delete: %w", err)
	}
	return nil
}

// DeleteAll removes every row in the table. Use with care.
func (r *Repository[T]) DeleteAll(ctx context.Context) error {
	var t T
	sqlStr, args := query.Delete(resolveTableName(t)).WithDialect(r.dialect).Build()
	if _, err := r.db.ExecContext(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("oca: delete all: %w", err)
	}
	return nil
}
