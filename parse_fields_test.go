package oca

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Example struct for testing schema
type TodoModel struct {
	ID        int       `db:"id,pk,auto"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at" schema:"default:now,notnull:true"`
	UpdatedAt time.Time `db:"updated_at" schema:"auto:true"`
	Skip      string    `db:"-"`
}

func TestParseFieldsWithSchema(t *testing.T) {
	todo := TodoModel{
		ID:    1,
		Title: "Test task",
	}

	fields := parseFields(todo)

	assert.Len(t, fields, 4, "Should parse 4 fields (skip ignored)")

	// Check ID field
	assert.Equal(t, "ID", fields[0].Name)
	assert.Equal(t, "id", fields[0].Column)
	assert.True(t, fields[0].IsPrimary)
	assert.True(t, fields[0].IsAuto)
	assert.Equal(t, 1, fields[0].Value)
	assert.Empty(t, fields[0].Schema)

	// Check Title field
	assert.Equal(t, "Title", fields[1].Name)
	assert.Equal(t, "title", fields[1].Column)
	assert.False(t, fields[1].IsPrimary)
	assert.Empty(t, fields[1].Schema)
	assert.Equal(t, "Test task", fields[1].Value)

	// Check CreatedAt field
	assert.Equal(t, "CreatedAt", fields[2].Name)
	assert.Equal(t, "created_at", fields[2].Column)
	assert.Equal(t, map[string]string{"default": "now", "notnull": "true"}, fields[2].Schema)

	// Check UpdatedAt field
	assert.Equal(t, "UpdatedAt", fields[3].Name)
	assert.Equal(t, "updated_at", fields[3].Column)
	assert.Equal(t, map[string]string{"auto": "true"}, fields[3].Schema)
}

func TestParseFieldsWithPointerAndSchema(t *testing.T) {
	todo := &TodoModel{ID: 2, Title: "Pointer task"}

	fields := parseFields(todo)
	assert.Len(t, fields, 4, "Should parse fields for pointer struct too")
	assert.Equal(t, 2, fields[0].Value)
}

func TestGetColumnNamesWithSchema(t *testing.T) {
	todo := TodoModel{}
	cols := getColumnNames(todo)

	assert.Equal(t, []string{"id", "title", "created_at", "updated_at"}, cols)
}

func TestCacheBehaviorWithSchema(t *testing.T) {
	todo := TodoModel{}
	typ := reflect.TypeOf(todo)

	// First call builds metadata
	metas1 := getStructMeta(typ)
	// Second call should hit cache
	metas2 := getStructMeta(typ)

	assert.Equal(t, metas1, metas2, "Cached metadata should be identical")
	assert.True(t, &metas1[0] != &metas2[0], "Should be different slices but equal content")
}
