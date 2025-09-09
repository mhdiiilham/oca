package query

import (
	"fmt"
	"strings"
)

// Condition represents a SQL WHERE condition expression and its bound arguments.
// Example: Condition{Expr: "age > ?", Args: []any{18}}
type Condition struct {
	Expr string // SQL expression with placeholders (e.g., "age > ?")
	Args []any  // Values to bind into placeholders
}

// C creates a new column reference that can be used to build conditions.
// Example: query.C("age").Gt(18)
func C(column string) Column {
	return Column{name: column}
}

// Column wraps a SQL column name and provides methods for building conditions.
type Column struct {
	name string
}

//
// --- Basic comparison operators ---
//

// Eq creates an equality condition: "column = ?".
func (c Column) Eq(val any) Condition {
	return Condition{Expr: fmt.Sprintf("%s = ?", c.name), Args: []any{val}}
}

// Neq creates a not equal condition: "column != ?".
func (c Column) Neq(val any) Condition {
	return Condition{Expr: fmt.Sprintf("%s != ?", c.name), Args: []any{val}}
}

// Gt creates a greater-than condition: "column > ?".
func (c Column) Gt(val any) Condition {
	return Condition{Expr: fmt.Sprintf("%s > ?", c.name), Args: []any{val}}
}

// Gte creates a greater-than-or-equal condition: "column >= ?".
func (c Column) Gte(val any) Condition {
	return Condition{Expr: fmt.Sprintf("%s >= ?", c.name), Args: []any{val}}
}

// Lt creates a less-than condition: "column < ?".
func (c Column) Lt(val any) Condition {
	return Condition{Expr: fmt.Sprintf("%s < ?", c.name), Args: []any{val}}
}

// Lte creates a less-than-or-equal condition: "column <= ?".
func (c Column) Lte(val any) Condition {
	return Condition{Expr: fmt.Sprintf("%s <= ?", c.name), Args: []any{val}}
}

//
// --- Null checks ---
//

// IsNull creates a condition: "column IS NULL".
func (c Column) IsNull() Condition {
	return Condition{Expr: fmt.Sprintf("%s IS NULL", c.name)}
}

// IsNotNull creates a condition: "column IS NOT NULL".
func (c Column) IsNotNull() Condition {
	return Condition{Expr: fmt.Sprintf("%s IS NOT NULL", c.name)}
}

//
// --- IN / NOT IN ---
//

// In creates a condition: "column IN (?, ?, ...)".
func (c Column) In(vals ...any) Condition {
	placeholders := strings.TrimRight(strings.Repeat("?,", len(vals)), ",")
	return Condition{
		Expr: fmt.Sprintf("%s IN (%s)", c.name, placeholders),
		Args: vals,
	}
}

// NotIn creates a condition: "column NOT IN (?, ?, ...)".
func (c Column) NotIn(vals ...any) Condition {
	placeholders := strings.TrimRight(strings.Repeat("?,", len(vals)), ",")
	return Condition{
		Expr: fmt.Sprintf("%s NOT IN (%s)", c.name, placeholders),
		Args: vals,
	}
}

//
// --- String matching ---
//

// Like creates a condition: "column LIKE ?".
// Example: query.C("name").Like("%john%")
func (c Column) Like(pattern string) Condition {
	return Condition{Expr: fmt.Sprintf("%s LIKE ?", c.name), Args: []any{pattern}}
}

// NotLike creates a condition: "column NOT LIKE ?".
func (c Column) NotLike(pattern string) Condition {
	return Condition{Expr: fmt.Sprintf("%s NOT LIKE ?", c.name), Args: []any{pattern}}
}

//
// --- Range matching ---
//

// Between creates a condition: "column BETWEEN ? AND ?".
func (c Column) Between(start, end any) Condition {
	return Condition{
		Expr: fmt.Sprintf("%s BETWEEN ? AND ?", c.name),
		Args: []any{start, end},
	}
}

// NotBetween creates a condition: "column NOT BETWEEN ? AND ?".
func (c Column) NotBetween(start, end any) Condition {
	return Condition{
		Expr: fmt.Sprintf("%s NOT BETWEEN ? AND ?", c.name),
		Args: []any{start, end},
	}
}

//
// --- Logical operators ---
//

// And combines multiple conditions with AND: "(cond1 AND cond2 ...)".
func And(conds ...Condition) Condition {
	exprs := make([]string, 0, len(conds))
	args := make([]any, 0)
	for _, cond := range conds {
		exprs = append(exprs, cond.Expr)
		args = append(args, cond.Args...)
	}
	return Condition{
		Expr: fmt.Sprintf("(%s)", strings.Join(exprs, " AND ")),
		Args: args,
	}
}

// Or combines multiple conditions with OR: "(cond1 OR cond2 ...)".
func Or(conds ...Condition) Condition {
	exprs := make([]string, 0, len(conds))
	args := make([]any, 0)
	for _, cond := range conds {
		exprs = append(exprs, cond.Expr)
		args = append(args, cond.Args...)
	}
	return Condition{
		Expr: fmt.Sprintf("(%s)", strings.Join(exprs, " OR ")),
		Args: args,
	}
}

// Not negates a condition: "NOT (cond)".
func Not(cond Condition) Condition {
	expr := cond.Expr

	// Remove outer parentheses if already wrapped
	if len(expr) > 1 && expr[0] == '(' && expr[len(expr)-1] == ')' {
		expr = expr[1 : len(expr)-1]
	}

	return Condition{
		Expr: fmt.Sprintf("NOT (%s)", expr),
		Args: cond.Args,
	}
}
