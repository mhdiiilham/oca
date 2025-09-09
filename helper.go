package oca

import (
	"fmt"
	"reflect"
	"strings"
)

// resolveTableName returns the table name either from Tabler or from struct name.
func resolveTableName[T any](entity T) string {
	if t, ok := any(entity).(Tabler); ok {
		return t.TableName()
	}

	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Fallback: use lowercase struct name
	return strings.ToLower(typ.Name()) + "s"
}

func buildScanTargets(val reflect.Value, fields []FieldMeta) ([]interface{}, error) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if !val.IsValid() {
		return nil, fmt.Errorf("buildScanTargets: invalid value")
	}

	targets := make([]interface{}, len(fields))
	for i, f := range fields {
		field := val.Field(f.Index)
		if !field.CanAddr() {
			return nil, fmt.Errorf("buildScanTargets: field %s is not addressable", f.Name)
		}
		targets[i] = field.Addr().Interface()
	}
	return targets, nil
}
