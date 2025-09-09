package query_test

import (
	"testing"

	"github.com/mhdiiilham/oca/query"
	"github.com/stretchr/testify/assert"
)

func TestInsertSingleRow(t *testing.T) {
	sql, args := query.InsertInto("users").
		Columns("id", "name").
		Values(1, "Alice").
		ToSQL()

	assert.Equal(t, "INSERT INTO users (id, name) VALUES (?, ?)", sql)
	assert.Equal(t, []interface{}{1, "Alice"}, args)
}

func TestInsertMultipleRows(t *testing.T) {
	sql, args := query.InsertInto("users").
		Columns("id", "name").
		Values(1, "Alice").
		Values(2, "Bob").
		ToSQL()

	assert.Equal(t, "INSERT INTO users (id, name) VALUES (?, ?), (?, ?)", sql)
	assert.Equal(t, []interface{}{1, "Alice", 2, "Bob"}, args)
}

func TestInsertWithReturning(t *testing.T) {
	sql, args := query.InsertInto("users").
		Columns("id", "name").
		Values(1, "Alice").
		Returning("id", "created_at").
		ToSQL()

	assert.Equal(t, "INSERT INTO users (id, name) VALUES (?, ?) RETURNING id, created_at", sql)
	assert.Equal(t, []interface{}{1, "Alice"}, args)
}

func TestInsertWithReturningID(t *testing.T) {
	sql, args := query.InsertInto("users").
		Columns("id", "name").
		Values(1, "Alice").
		ReturningID().
		ToSQL()

	assert.Equal(t, "INSERT INTO users (id, name) VALUES (?, ?) RETURNING id", sql)
	assert.Equal(t, []interface{}{1, "Alice"}, args)
}

func TestInsertWithRawSQL(t *testing.T) {
	sql, args := query.InsertInto("users").
		Columns("name", "created_at").
		Values("Alice", query.Raw("NOW()")).
		ToSQL()

	assert.Equal(t, "INSERT INTO users (name, created_at) VALUES (?, NOW())", sql)
	assert.Equal(t, []interface{}{"Alice"}, args)
}

func TestInsertEmptyBuilder(t *testing.T) {
	sql, args := query.InsertInto("").ToSQL()
	assert.Equal(t, "", sql)
	assert.Nil(t, args)

	sql, args = query.InsertInto("users").ToSQL()
	assert.Equal(t, "", sql)
	assert.Nil(t, args)
}
