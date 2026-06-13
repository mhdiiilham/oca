package oca

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

// scanRows maps sql.Rows into a slice of T using cached FieldMeta for column → field mapping.
func scanRows[T any](rows *sql.Rows) ([]T, error) {
	defer rows.Close()

	var result []T
	var entity T
	columns := getColumnNames(entity)
	colMap := buildColumnMap(entity)

	for rows.Next() {
		val := reflect.New(reflect.TypeOf(entity)).Elem()
		scanTargets := make([]any, len(columns))

		for i, col := range columns {
			index, ok := colMap[col]
			if !ok {
				return nil, fmt.Errorf("oca: scanRows: cannot map column %s", col)
			}
			scanTargets[i] = val.FieldByIndex([]int{index}).Addr().Interface()
		}

		if err := rows.Scan(scanTargets...); err != nil {
			return nil, err
		}

		result = append(result, val.Interface().(T))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// scanRowSingle maps a single sql.Row into T.
func scanRowSingle[T any](row *sql.Row, columns []string, colMap map[string]int) (*T, error) {
	var entity T
	val := reflect.New(reflect.TypeOf(entity)).Elem()
	scanTargets := make([]any, len(columns))

	for i, col := range columns {
		index, ok := colMap[col]
		if !ok {
			return nil, fmt.Errorf("oca: scanRowSingle: cannot map column %s", col)
		}
		scanTargets[i] = val.FieldByIndex([]int{index}).Addr().Interface()
	}

	if err := row.Scan(scanTargets...); err != nil {
		return nil, err
	}

	t := val.Interface().(T)
	return &t, nil
}

// scanAutoFields scans auto-generated values (e.g. RETURNING id) back onto entity.
func scanAutoFields[T any](ctx context.Context, db executor, sqlStr string, sqlArgs []any, entity *T, autoFields []FieldMeta) error {
	row := db.QueryRowContext(ctx, sqlStr, sqlArgs...)
	val := reflect.ValueOf(entity).Elem()
	scanTargets := make([]any, len(autoFields))

	for i, f := range autoFields {
		field := val.Field(f.Index)
		if !field.CanAddr() {
			return fmt.Errorf("oca: scanAutoFields: field %s is not addressable", f.Name)
		}
		scanTargets[i] = field.Addr().Interface()
	}

	return row.Scan(scanTargets...)
}

// buildColumnMap maps struct column names to their field indices.
func buildColumnMap(entity any) map[string]int {
	colMap := make(map[string]int)
	metas := getStructMeta(reflect.TypeOf(entity))
	for _, m := range metas {
		colMap[m.Column] = m.Index
	}
	return colMap
}
