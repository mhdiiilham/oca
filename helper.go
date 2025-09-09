package oca

import (
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
