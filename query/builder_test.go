package query_test

import (
	"testing"

	"github.com/mhdiiilham/oca/query"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_AllDialects(t *testing.T) {
	dialects := []struct {
		name    string
		dialect query.Dialect
	}{
		{"MySQL", query.MySQLDialect{}},
		{"MariaDB", query.MySQLDialect{}},
		{"Postgres", query.PostgresDialect{}},
	}

	for _, d := range dialects {
		d := d
		t.Run(d.name, func(t *testing.T) {
			query.SetDialect(d.dialect)

			// Simple SELECT
			sql, args := query.From("users").Select("id", "name").Build()
			assert.Equal(t, "SELECT id, name FROM users", sql)
			assert.Empty(t, args)

			// Default SELECT *
			sql, args = query.From("products").Build()
			assert.Equal(t, "SELECT * FROM products", sql)
			assert.Empty(t, args)

			// WHERE Eq
			sql, args = query.From("users").Select("id").Where(query.C("age").Gt(18)).Build()
			if d.name == "Postgres" {
				assert.Equal(t, "SELECT id FROM users WHERE age > $1", sql)
			} else {
				assert.Equal(t, "SELECT id FROM users WHERE age > ?", sql)
			}
			assert.Equal(t, []any{18}, args)

			// Multiple WHERE
			sql, args = query.From("users").
				Select("id").
				Where(query.C("age").Gt(18)).
				Where(query.C("status").Eq("active")).
				Build()
			if d.name == "Postgres" {
				assert.Equal(t, "SELECT id FROM users WHERE age > $1 AND status = $2", sql)
			} else {
				assert.Equal(t, "SELECT id FROM users WHERE age > ? AND status = ?", sql)
			}
			assert.Equal(t, []any{18, "active"}, args)

			// WHERE IN
			sql, args = query.From("users").
				Select("id").
				Where(query.C("id").In(1, 2, 3)).
				Build()
			if d.name == "Postgres" {
				assert.Equal(t, "SELECT id FROM users WHERE id IN ($1,$2,$3)", sql)
			} else {
				assert.Equal(t, "SELECT id FROM users WHERE id IN (?,?,?)", sql)
			}
			assert.Equal(t, []any{1, 2, 3}, args)

			// ORDER BY + LIMIT + OFFSET
			sql, args = query.From("users").
				Select("id").
				OrderBy("created_at DESC").
				Limit(10).
				Offset(5).
				Build()
			if d.name == "Postgres" {
				assert.Equal(t, "SELECT id FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2", sql)
			} else {
				assert.Equal(t, "SELECT id FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?", sql)
			}
			assert.Equal(t, []any{10, 5}, args)
		})
	}
}
