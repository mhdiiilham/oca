package query_test

import (
	"testing"

	"github.com/mhdiiilham/oca/query"
	"github.com/stretchr/testify/assert"
)

func TestInnerJoin(t *testing.T) {
	sql, args := query.From("users").
		Select("users.id", "profiles.bio").
		Join("profiles", "users.id = profiles.user_id").
		Build()

	assert.Equal(t, "SELECT users.id, profiles.bio FROM users INNER JOIN profiles ON users.id = profiles.user_id", sql)
	assert.Empty(t, args)
}

func TestLeftJoin(t *testing.T) {
	sql, args := query.From("orders").
		Select("orders.id", "customers.name").
		LeftJoin("customers", "orders.customer_id = customers.id").
		Build()

	assert.Equal(t, "SELECT orders.id, customers.name FROM orders LEFT JOIN customers ON orders.customer_id = customers.id", sql)
	assert.Empty(t, args)
}

func TestRightJoin(t *testing.T) {
	sql, args := query.From("a").
		Select("a.x", "b.y").
		RightJoin("b", "a.id = b.a_id").
		Build()

	assert.Equal(t, "SELECT a.x, b.y FROM a RIGHT JOIN b ON a.id = b.a_id", sql)
	assert.Empty(t, args)
}

func TestFullJoin(t *testing.T) {
	sql, args := query.From("products").
		Select("products.id", "categories.name").
		FullJoin("categories", "products.category_id = categories.id").
		Build()

	assert.Equal(t, "SELECT products.id, categories.name FROM products FULL JOIN categories ON products.category_id = categories.id", sql)
	assert.Empty(t, args)
}
