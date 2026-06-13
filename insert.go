package oca

import (
	"context"
	"fmt"
	"reflect"

	"github.com/mhdiiilham/oca/query"
)

// Insert inserts one entity into the database.
// Auto-increment fields (db:"...,auto") and schema defaults (schema:"default:now()")
// are handled automatically; auto fields are populated on entity after insert.
func (r *Repository[T]) Insert(ctx context.Context, entity *T) error {
	cols, args, autoFields, err := prepareInsertFields(entity)
	if err != nil {
		return err
	}

	var t T
	builder := query.InsertInto(resolveTableName(t)).WithDialect(r.dialect).
		Columns(cols...).
		Values(args...)

	if len(autoFields) > 0 {
		autoCols := make([]string, len(autoFields))
		for i, f := range autoFields {
			autoCols[i] = f.Column
		}
		builder = builder.Returning(autoCols...)
	}

	sqlStr, sqlArgs := builder.ToSQL()
	if sqlStr == "" {
		return fmt.Errorf("oca: insert: nothing to insert")
	}

	if len(autoFields) > 0 {
		return r.scanAutoFields(ctx, entity, sqlStr, sqlArgs, autoFields)
	}

	if _, err := r.db.ExecContext(ctx, sqlStr, sqlArgs...); err != nil {
		return fmt.Errorf("oca: insert: %w", err)
	}
	return nil
}

// InsertMany inserts multiple entities in a single query.
// Auto-increment fields are not populated back; use Insert for that.
// Nil pointer fields are sent as NULL.
func (r *Repository[T]) InsertMany(ctx context.Context, entities []T) error {
	if len(entities) == 0 {
		return nil
	}

	var t T
	tableName := resolveTableName(t)

	// Derive the fixed column layout from the first entity.
	sampleFields := parseFields(&entities[0])

	type colInfo struct {
		col       string
		fieldIdx  int
		isDefault bool // schema:"default:now()"
	}

	var layout []colInfo
	for _, f := range sampleFields {
		if f.IsAuto {
			continue
		}
		layout = append(layout, colInfo{
			col:       f.Column,
			fieldIdx:  f.Index,
			isDefault: f.Schema["default"] == "now()",
		})
	}

	colNames := make([]string, len(layout))
	for i, c := range layout {
		colNames[i] = c.col
	}

	builder := query.InsertInto(tableName).WithDialect(r.dialect).Columns(colNames...)

	for i := range entities {
		ev := reflect.ValueOf(&entities[i]).Elem()
		row := make([]any, len(layout))
		for j, c := range layout {
			if c.isDefault {
				row[j] = query.Raw("NOW()")
				continue
			}
			fv := ev.Field(c.fieldIdx)
			if fv.Kind() == reflect.Pointer {
				if fv.IsNil() {
					row[j] = nil
				} else {
					row[j] = fv.Elem().Interface()
				}
			} else {
				row[j] = fv.Interface()
			}
		}
		builder.Values(row...)
	}

	sqlStr, sqlArgs := builder.ToSQL()
	if sqlStr == "" {
		return fmt.Errorf("oca: insert many: empty query")
	}

	if _, err := r.db.ExecContext(ctx, sqlStr, sqlArgs...); err != nil {
		return fmt.Errorf("oca: insert many: %w", err)
	}
	return nil
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

func (r *Repository[T]) scanAutoFields(ctx context.Context, entity *T, sqlStr string, sqlArgs []any, autoFields []FieldMeta) error {
	if entity == nil {
		return fmt.Errorf("oca: scanAutoFields: entity cannot be nil")
	}

	row := r.db.QueryRowContext(ctx, sqlStr, sqlArgs...)
	val := reflect.ValueOf(entity)
	targets, err := buildScanTargets(val, autoFields)
	if err != nil {
		return err
	}

	return row.Scan(targets...)
}
