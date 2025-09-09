package query_test

import (
	"testing"

	"github.com/mhdiiilham/oca/query"
	"github.com/stretchr/testify/assert"
)

func TestInsertBuilder_Basic(t *testing.T) {
	q := query.InsertInto("users").
		Columns("name", "email").
		Values("Alice", "alice@example.com")

	sql, args := q.ToSQL()

	expectedSQL := "INSERT INTO users (name, email) VALUES (?, ?)"
	expectedArgs := []interface{}{"Alice", "alice@example.com"}
	assert.Equal(t, expectedSQL, sql)
	assert.Equal(t, expectedArgs, args)

	for i, v := range expectedArgs {
		assert.Equal(t, v, args[i])
	}
}

func TestInsertBuilder_Returning(t *testing.T) {
	q := query.InsertInto("users").
		Columns("name", "email").
		Values("Bob", "bob@example.com").
		Returning("id", "created_at")

	sql, args := q.ToSQL()

	expectedSQL := "INSERT INTO users (name, email) VALUES (?, ?) RETURNING id, created_at"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{"Bob", "bob@example.com"}
	assert.Equal(t, expectedArgs, args)

	for i, v := range expectedArgs {
		assert.Equal(t, v, args[i])
	}
}

func TestInsertBuilder_ReturningID(t *testing.T) {
	q := query.InsertInto("users").
		Columns("username").
		Values("charlie").
		ReturningID()

	sql, args := q.ToSQL()

	expectedSQL := "INSERT INTO users (username) VALUES (?) RETURNING id"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{"charlie"}
	assert.Equal(t, expectedArgs, args)

	for i, v := range expectedArgs {
		assert.Equal(t, v, args[i])
	}
}

// New test for RawSQL support (e.g., default: now())
func TestInsertBuilder_RawSQLValue(t *testing.T) {
	q := query.InsertInto("todos").
		Columns("title", "created_at").
		Values("Task 1", query.Raw("NOW()"))

	sql, args := q.ToSQL()

	expectedSQL := "INSERT INTO todos (title, created_at) VALUES (?, NOW())"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{"Task 1"} // only the Go value
	assert.Equal(t, expectedArgs, args)
}
