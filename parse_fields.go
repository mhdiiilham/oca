package oca

import (
	"reflect"
	"strings"
	"sync"
)

// FieldMeta holds metadata about a struct field for ORM mapping.
type FieldMeta struct {
	Name      string // Go struct field name
	Column    string // DB column name from tag
	Index     int    // field index in struct
	IsPrimary bool   // true if marked as primary key
	IsAuto    bool   // true if marked as auto increment
}

var (
	fieldCache sync.Map // map[reflect.Type][]FieldMeta
)

// getStructMeta returns cached field metadata for a struct type.
// It stores one canonical slice in the cache but always returns a copy
// so callers cannot mutate the cached metadata.
func getStructMeta(t reflect.Type) []FieldMeta {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// check cache
	if cached, ok := fieldCache.Load(t); ok {
		orig := cached.([]FieldMeta)
		// return a copy to avoid exposing the cached backing array to callers
		out := make([]FieldMeta, len(orig))
		copy(out, orig)
		return out
	}

	// build metadata
	var fields []FieldMeta
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		tag := sf.Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}

		parts := strings.Split(tag, ",")
		meta := FieldMeta{
			Name:   sf.Name,
			Column: parts[0],
			Index:  i,
		}

		for _, opt := range parts[1:] {
			switch strings.TrimSpace(opt) {
			case "pk":
				meta.IsPrimary = true
			case "auto":
				meta.IsAuto = true
			}
		}
		fields = append(fields, meta)
	}

	// store canonical slice in cache
	fieldCache.Store(t, fields)

	// return a copy to callers
	out := make([]FieldMeta, len(fields))
	copy(out, fields)
	return out
}

// parseFields extracts field metadata and current values.
// This function is cheap after first call thanks to caching.
func parseFields(entity any) []struct {
	FieldMeta
	Value interface{}
} {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	metas := getStructMeta(val.Type())
	result := make([]struct {
		FieldMeta
		Value interface{}
	}, len(metas))

	for i, m := range metas {
		result[i] = struct {
			FieldMeta
			Value interface{}
		}{
			FieldMeta: m,
			Value:     val.Field(m.Index).Interface(),
		}
	}

	return result
}

// getColumnNames returns only the DB column names.
func getColumnNames(entity any) []string {
	metas := getStructMeta(reflect.TypeOf(entity))
	cols := make([]string, len(metas))
	for i, m := range metas {
		cols[i] = m.Column
	}
	return cols
}
