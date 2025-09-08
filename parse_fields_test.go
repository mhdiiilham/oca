package oca

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Example struct for testing
type User struct {
	ID    int    `db:"id,pk,auto"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Skip  string `db:"-"` // ignored
}

func TestParseFields(t *testing.T) {
	u := User{ID: 1, Name: "Alice", Email: "alice@example.com"}

	fields := parseFields(u)

	assert.Len(t, fields, 3, "Should only parse 3 fields (skip ignored)")

	// Check field metadata
	assert.Equal(t, "ID", fields[0].Name)
	assert.Equal(t, "id", fields[0].Column)
	assert.True(t, fields[0].IsPrimary)
	assert.True(t, fields[0].IsAuto)
	assert.Equal(t, 1, fields[0].Value)

	assert.Equal(t, "Name", fields[1].Name)
	assert.Equal(t, "name", fields[1].Column)
	assert.False(t, fields[1].IsPrimary)
	assert.Equal(t, "Alice", fields[1].Value)

	assert.Equal(t, "Email", fields[2].Name)
	assert.Equal(t, "email", fields[2].Column)
	assert.Equal(t, "alice@example.com", fields[2].Value)
}

func TestParseFieldsWithPointer(t *testing.T) {
	u := &User{ID: 2, Name: "Bob", Email: "bob@example.com"}

	fields := parseFields(u)
	assert.Len(t, fields, 3, "Should parse fields for pointer struct too")
	assert.Equal(t, 2, fields[0].Value)
}

func TestGetColumnNames(t *testing.T) {
	u := User{}
	cols := getColumnNames(u)

	assert.Equal(t, []string{"id", "name", "email"}, cols)
}

func TestCacheBehavior(t *testing.T) {
	u := User{}
	typ := reflect.TypeOf(u)

	// First call builds metadata
	metas1 := getStructMeta(typ)
	// Second call should hit cache
	metas2 := getStructMeta(typ)

	assert.Equal(t, metas1, metas2, "Cached metadata should be identical")
	assert.True(t, &metas1[0] != &metas2[0], "Should be different slices but equal content")
}
