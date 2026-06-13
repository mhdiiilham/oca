package query_test

import (
	"testing"

	"github.com/mhdiiilham/oca/query"
	"github.com/stretchr/testify/assert"
)

func TestUpdateBuilder_MySQL(t *testing.T) {
	query.SetDialect(query.MySQLDialect{})

	q := query.Update("users").
		Set("name", "Alice").
		Where(query.C("id").Eq(1))

	sql, args := q.Build()
	assert.Equal(t, "UPDATE users SET name = ? WHERE id = ?", sql)
	assert.Equal(t, []any{"Alice", 1}, args)
}

func TestUpdateBuilder_Postgres(t *testing.T) {
	query.SetDialect(query.PostgresDialect{})

	q := query.Update("users").
		Set("name", "Bob").
		Set("email", "bob@example.com").
		Where(query.C("id").Eq(42))

	sql, args := q.Build()
	assert.Equal(t, "UPDATE users SET name = $1, email = $2 WHERE id = $3", sql)
	assert.Equal(t, []any{"Bob", "bob@example.com", 42}, args)
}

func TestUpdateBuilder_MariaDB(t *testing.T) {
	query.SetDialect(query.MariaDBDialect{})

	q := query.Update("products").
		Set("price", 9.99).
		Where(query.C("sku").Eq("ABC-001"))

	sql, args := q.Build()
	assert.Equal(t, "UPDATE products SET price = ? WHERE sku = ?", sql)
	assert.Equal(t, []any{9.99, "ABC-001"}, args)
}

func TestUpdateBuilder_MultipleSet(t *testing.T) {
	query.SetDialect(query.MySQLDialect{})

	q := query.Update("orders").
		Set("status", "shipped").
		Set("tracking_number", "TRACK123").
		Set("updated_at", "2026-06-13").
		Where(query.C("id").Eq(100))

	sql, args := q.Build()
	assert.Equal(t, "UPDATE orders SET status = ?, tracking_number = ?, updated_at = ? WHERE id = ?", sql)
	assert.Equal(t, []any{"shipped", "TRACK123", "2026-06-13", 100}, args)
}

func TestUpdateBuilder_MultipleWhere(t *testing.T) {
	query.SetDialect(query.MySQLDialect{})

	q := query.Update("sessions").
		Set("active", false).
		Where(
			query.C("user_id").Eq(7),
			query.C("expired").Eq(true),
		)

	sql, args := q.Build()
	assert.Equal(t, "UPDATE sessions SET active = ? WHERE user_id = ? AND expired = ?", sql)
	assert.Equal(t, []any{false, 7, true}, args)
}

func TestUpdateBuilder_NoWhere(t *testing.T) {
	query.SetDialect(query.MySQLDialect{})

	q := query.Update("settings").
		Set("maintenance_mode", true)

	sql, args := q.Build()
	assert.Equal(t, "UPDATE settings SET maintenance_mode = ?", sql)
	assert.Equal(t, []any{true}, args)
}
