package query_test

import (
	"testing"

	"github.com/mhdiiilham/oca/query"
	"github.com/stretchr/testify/assert"
)

func TestColumnOperators(t *testing.T) {
	c := query.C("age")

	// --- Comparison ---
	cond := c.Eq(18)
	assert.Equal(t, "age = ?", cond.Expr)
	assert.Equal(t, []any{18}, cond.Args)

	cond = c.Neq(20)
	assert.Equal(t, "age != ?", cond.Expr)
	assert.Equal(t, []any{20}, cond.Args)

	cond = c.Gt(10)
	assert.Equal(t, "age > ?", cond.Expr)
	assert.Equal(t, []any{10}, cond.Args)

	cond = c.Gte(15)
	assert.Equal(t, "age >= ?", cond.Expr)
	assert.Equal(t, []any{15}, cond.Args)

	cond = c.Lt(50)
	assert.Equal(t, "age < ?", cond.Expr)
	assert.Equal(t, []any{50}, cond.Args)

	cond = c.Lte(40)
	assert.Equal(t, "age <= ?", cond.Expr)
	assert.Equal(t, []any{40}, cond.Args)

	// --- Null checks ---
	cond = c.IsNull()
	assert.Equal(t, "age IS NULL", cond.Expr)
	assert.Empty(t, cond.Args)

	cond = c.IsNotNull()
	assert.Equal(t, "age IS NOT NULL", cond.Expr)
	assert.Empty(t, cond.Args)

	// --- IN / NOT IN ---
	cond = c.In(1, 2, 3)
	assert.Equal(t, "age IN (?,?,?)", cond.Expr)
	assert.Equal(t, []any{1, 2, 3}, cond.Args)

	cond = c.NotIn(4, 5, 6)
	assert.Equal(t, "age NOT IN (?,?,?)", cond.Expr)
	assert.Equal(t, []any{4, 5, 6}, cond.Args)

	// --- LIKE / NOT LIKE ---
	cond = c.Like("%john%")
	assert.Equal(t, "age LIKE ?", cond.Expr)
	assert.Equal(t, []any{"%john%"}, cond.Args)

	cond = c.NotLike("%doe%")
	assert.Equal(t, "age NOT LIKE ?", cond.Expr)
	assert.Equal(t, []any{"%doe%"}, cond.Args)

	// --- BETWEEN / NOT BETWEEN ---
	cond = c.Between(10, 20)
	assert.Equal(t, "age BETWEEN ? AND ?", cond.Expr)
	assert.Equal(t, []any{10, 20}, cond.Args)

	cond = c.NotBetween(30, 40)
	assert.Equal(t, "age NOT BETWEEN ? AND ?", cond.Expr)
	assert.Equal(t, []any{30, 40}, cond.Args)
}

func TestLogicalOperators(t *testing.T) {
	c1 := query.C("age").Gt(18)
	c2 := query.C("status").Eq("active")
	c3 := query.C("score").Lt(100)

	// --- AND ---
	andCond := query.And(c1, c2, c3)
	assert.Equal(t, "(age > ? AND status = ? AND score < ?)", andCond.Expr)
	assert.Equal(t, []any{18, "active", 100}, andCond.Args)

	// --- OR ---
	orCond := query.Or(c1, c2, c3)
	assert.Equal(t, "(age > ? OR status = ? OR score < ?)", orCond.Expr)
	assert.Equal(t, []any{18, "active", 100}, orCond.Args)

	// --- NOT ---
	notCond := query.Not(c1)
	assert.Equal(t, "NOT (age > ?)", notCond.Expr)
	assert.Equal(t, []any{18}, notCond.Args)
}

func TestNestedLogicalOperators(t *testing.T) {
	c1 := query.C("age").Gt(18)
	c2 := query.C("status").Eq("active")
	c3 := query.C("score").Lt(100)

	// AND inside OR
	cond := query.Or(c1, query.And(c2, c3))
	assert.Equal(t, "(age > ? OR (status = ? AND score < ?))", cond.Expr)
	assert.Equal(t, []any{18, "active", 100}, cond.Args)

	// NOT with AND
	cond2 := query.Not(query.And(c1, c2))
	assert.Equal(t, "NOT (age > ? AND status = ?)", cond2.Expr)
	assert.Equal(t, []any{18, "active"}, cond2.Args)
}
