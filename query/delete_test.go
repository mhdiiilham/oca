package query_test

import (
	"testing"

	"github.com/mhdiiilham/oca/query"
	"github.com/stretchr/testify/assert"
)

func TestDeleteBuilder_MySQL(t *testing.T) {
	query.SetDialect(query.MySQLDialect{})

	q := query.Delete("users").
		Where(query.C("id").Eq(10))

	sql, args := q.Build()
	expectedSQL := "DELETE FROM users WHERE id = ?"
	expectedArgs := []any{10}

	assert.Equal(t, expectedSQL, sql, "got %q, want %q", sql, expectedSQL)
	assert.Equal(t, len(expectedArgs), len(args), "got args %v, want %v", args, expectedArgs)
	for i := range expectedArgs {
		assert.Equal(t, expectedArgs[i], args[i], "got args %v, want %v", args[i], expectedArgs[i])
	}
}

func TestDeleteBuilder_Postgres(t *testing.T) {
	query.SetDialect(query.PostgresDialect{})

	q := query.Delete("orders").
		Where(
			query.C("status").Eq("cancelled"),
			query.C("id").Eq(99),
		)

	sql, args := q.Build()
	expectedSQL := "DELETE FROM orders WHERE status = $1 AND id = $2"
	expectedArgs := []any{"cancelled", 99}

	assert.Equal(t, expectedSQL, sql, "got %q, want %q", sql, expectedSQL)
	assert.Equal(t, len(expectedArgs), len(args), "got args %v, want %v", len(args), len(expectedArgs))

	for i := range args {
		assert.Equal(t, expectedArgs[i], args[i], "arg %d: got %v, want %v", i, args[i], expectedArgs[i])
	}
}

func TestDeleteBuilder_NoWhere(t *testing.T) {
	query.SetDialect(query.MySQLDialect{})

	q := query.Delete("sessions")

	sql, args := q.Build()
	expectedSQL := "DELETE FROM sessions"
	assert.Equal(t, expectedSQL, sql, "got %q, want %q", sql, expectedSQL)
	assert.Empty(t, args, "expected no args, got %v", args)
}
