package oca

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

// scanRows maps sql.Rows into a slice of structs of type T.
// It uses reflection and cached FieldMeta for fast column -> field mapping.
func scanRows[T any](rows *sql.Rows) ([]T, error) {
	defer rows.Close()

	var result []T
	var entity T
	columns := getColumnNames(entity)
	colMap := buildColumnMap(entity) // map column name -> []int

	for rows.Next() {
		val := reflect.New(reflect.TypeOf(entity)).Elem()
		scanTargets := make([]interface{}, len(columns))

		for i, col := range columns {
			index, ok := colMap[col]
			if !ok {
				return nil, fmt.Errorf("scanRows: cannot map column %s", col)
			}
			// Use slice of int for FieldByIndex
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

// scanRowSingle maps sql.Row (single row) into a struct of type T.
// Useful for FindOne queries or returning auto-generated IDs.
func scanRowSingle[T any](row *sql.Row, columns []string, colMap map[string]int) (*T, error) {
	var entity T
	val := reflect.New(reflect.TypeOf(entity)).Elem()
	scanTargets := make([]interface{}, len(columns))

	for i, col := range columns {
		index, ok := colMap[col]
		if !ok {
			return nil, fmt.Errorf("scanRowSingle: cannot map column %s", col)
		}
		scanTargets[i] = val.FieldByIndex([]int{index}).Addr().Interface()
	}

	if err := row.Scan(scanTargets...); err != nil {
		return nil, err
	}

	t := val.Interface().(T)
	return &t, nil
}

// scanAutoFields scans values from sql.Row into specified auto-increment fields.
// Useful after INSERT with RETURNING or lastInsertId.
func scanAutoFields[T any](ctx context.Context, db *sql.DB, sqlStr string, sqlArgs []interface{}, entity *T, autoFields []FieldMeta) error {
	row := db.QueryRowContext(ctx, sqlStr, sqlArgs...)
	val := reflect.ValueOf(entity).Elem()
	scanTargets := make([]interface{}, len(autoFields))

	for i, f := range autoFields {
		field := val.Field(f.Index)
		if !field.CanAddr() {
			return &reflect.ValueError{Method: "scanAutoFields", Kind: field.Kind()}
		}
		scanTargets[i] = field.Addr().Interface()
	}

	return row.Scan(scanTargets...)
}

// buildColumnMap maps struct columns to their field indices.
func buildColumnMap(entity any) map[string]int {
	colMap := make(map[string]int)
	metas := getStructMeta(reflect.TypeOf(entity))
	for _, m := range metas {
		colMap[m.Column] = m.Index
	}
	return colMap
}
