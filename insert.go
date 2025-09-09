package oca

import (
	"context"
	"fmt"
	"reflect"

	"github.com/mhdiiilham/oca/query"
)

// Insert inserts the entity into the database.
// It automatically handles auto-increment fields and schema defaults (e.g., default:now()).
// Multiple auto fields are supported and populated after insertion.
func (r *Repository[T]) Insert(ctx context.Context, entity *T) error {
	cols, args, autoFields, err := prepareInsertFields(entity)
	if err != nil {
		return err
	}

	sqlStr, sqlArgs := buildInsertQuery(entity, cols, args, autoFields)

	if len(autoFields) > 0 {
		return r.scanAutoFields(ctx, entity, sqlStr, sqlArgs, autoFields)
	}

	_, err = r.db.ExecContext(ctx, sqlStr, sqlArgs...)
	return err
}

func prepareInsertFields[T any](entity *T) (cols []string, args []any, autoFields []FieldMeta, err error) {
	fields := parseFields(entity)
	for _, f := range fields {
		switch {
		case f.IsAuto:
			autoFields = append(autoFields, f.FieldMeta)
		case f.Schema["default"] == "now()":
			cols = append(cols, f.Column)
			args = append(args, query.Raw("NOW()"))
		default:
			cols = append(cols, f.Column)
			args = append(args, f.Value)
		}
	}
	return
}

func buildInsertQuery[T any](entity *T, cols []string, args []any, autoFields []FieldMeta) (string, []interface{}) {
	var t T
	builder := query.InsertInto(resolveTableName(t)).
		Columns(cols...).
		Values(args...)

	if len(autoFields) > 0 {
		autoCols := make([]string, len(autoFields))
		for i, f := range autoFields {
			autoCols[i] = f.Column
		}
		builder = builder.Returning(autoCols...)
	}

	return builder.ToSQL()
}

func (r *Repository[T]) scanAutoFields(ctx context.Context, entity *T, sqlStr string, sqlArgs []interface{}, autoFields []FieldMeta) error {
	if entity == nil {
		return fmt.Errorf("scanAutoFields: entity cannot be nil")
	}

	row := r.db.QueryRowContext(ctx, sqlStr, sqlArgs...)
	val := reflect.ValueOf(entity)
	targets, err := buildScanTargets(val, autoFields)
	if err != nil {
		return err
	}

	return row.Scan(targets...)
}
