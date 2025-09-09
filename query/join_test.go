package query_test

import (
	"testing"

	"github.com/mhdiiilham/oca/query"
	"github.com/stretchr/testify/assert"
)

func TestJoinClauses(t *testing.T) {
	// Ensure default dialect is PostgreSQL
	query.SetDialect(query.PostgresDialect{})

	// INNER JOIN
	sql, _ := query.From("users").
		Select("users.id", "profiles.bio").
		Join("profiles", "users.id = profiles.user_id").
		Build()
	assert.Equal(t,
		"SELECT users.id, profiles.bio FROM users INNER JOIN profiles ON users.id = profiles.user_id",
		sql)

	// LEFT JOIN
	sql, _ = query.From("users").
		Select("users.id", "orders.amount").
		LeftJoin("orders", "users.id = orders.user_id").
		Build()
	assert.Equal(t,
		"SELECT users.id, orders.amount FROM users LEFT JOIN orders ON users.id = orders.user_id",
		sql)

	// RIGHT JOIN
	sql, _ = query.From("users").
		Select("users.id", "payments.status").
		RightJoin("payments", "users.id = payments.user_id").
		Build()
	assert.Equal(t,
		"SELECT users.id, payments.status FROM users RIGHT JOIN payments ON users.id = payments.user_id",
		sql)

	// FULL JOIN
	sql, _ = query.From("users").
		Select("users.id", "logs.action").
		FullJoin("logs", "users.id = logs.user_id").
		Build()
	assert.Equal(t,
		"SELECT users.id, logs.action FROM users FULL JOIN logs ON users.id = logs.user_id",
		sql)

	// JOIN + WHERE with PostgreSQL placeholders ($1)
	sql, args := query.From("users").
		Select("users.id", "orders.amount").
		Join("orders", "users.id = orders.user_id").
		Where(query.C("orders.amount").Gt(100)).
		Build()
	assert.Equal(t,
		"SELECT users.id, orders.amount FROM users INNER JOIN orders ON users.id = orders.user_id WHERE orders.amount > $1",
		sql)
	assert.Equal(t, []any{100}, args)
}
