package query_test

import (
	"testing"

	"github.com/mhdiiilham/oca/query"
	"github.com/stretchr/testify/assert"
)

func TestSimpleSelect(t *testing.T) {
	sql, args := query.From("users").Select("id", "name").Build()

	assert.Equal(t, "SELECT id, name FROM users", sql)
	assert.Empty(t, args)
}

func TestSelectWithWhere(t *testing.T) {
	sql, args := query.From("users").
		Select("id", "name").
		Where("age > ?", 18).
		Build()

	assert.Equal(t, "SELECT id, name FROM users WHERE age > ?", sql)
	assert.Equal(t, []any{18}, args)
}

func TestMultipleWhereClauses(t *testing.T) {
	sql, args := query.From("users").
		Select("id").
		Where("age > ?", 18).
		Where("status = ?", "active").
		Build()

	assert.Equal(t, "SELECT id FROM users WHERE age > ? AND status = ?", sql)
	assert.Equal(t, []any{18, "active"}, args)
}

func TestOrderByLimitOffset(t *testing.T) {
	sql, args := query.From("users").
		Select("id", "email").
		OrderBy("created_at DESC").
		Limit(10).
		Offset(5).
		Build()

	assert.Equal(t, "SELECT id, email FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?", sql)
	assert.Equal(t, []any{10, 5}, args)
}

func TestDefaultSelectAll(t *testing.T) {
	sql, args := query.From("products").Build()

	assert.Equal(t, "SELECT * FROM products", sql)
	assert.Empty(t, args)
}
