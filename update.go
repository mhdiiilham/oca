package oca

import (
	"context"
	"fmt"
	"reflect"

	"github.com/mhdiiilham/oca/query"
)

// Update writes all non-auto, non-pk fields of entity back to the database.
// Pointer fields that are nil are skipped (partial update — that column is untouched).
// Without Where() opts the primary key value is used as the WHERE condition.
func (r *Repository[T]) Update(ctx context.Context, entity *T, opts ...FilterOptions) error {
	fields := parseFields(entity)

	var t T
	builder := query.Update(resolveTableName(t)).WithDialect(r.dialect)

	filter := defaultFilter()
	for _, opt := range opts {
		if opt != nil {
			opt(&filter)
		}
	}

	var pkColumn string
	var pkValue any

	for _, f := range fields {
		if f.IsPrimary {
			pkColumn = f.Column
			pkValue = f.Value
			continue
		}
		if f.IsAuto {
			continue
		}
		// Skip nil pointer fields — caller did not set this field, leave it unchanged.
		fv := reflect.ValueOf(entity).Elem().Field(f.Index)
		if fv.Kind() == reflect.Pointer && fv.IsNil() {
			continue
		}
		builder.Set(f.Column, f.Value)
	}

	if len(filter.Where) > 0 {
		builder.Where(filter.Where...)
	} else if pkColumn != "" {
		builder.Where(query.C(pkColumn).Eq(pkValue))
	}

	sqlStr, args := builder.Build()
	if _, err := r.db.ExecContext(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("oca: update: %w", err)
	}
	return nil
}
