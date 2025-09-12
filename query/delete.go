package query

import (
	"strings"
)

// DeleteBuilder builds SQL DELETE queries with dialect support.
type DeleteBuilder struct {
	table            string
	where            []Condition
	args             []any
	placeholderIndex int
}

// Delete creates a new DeleteBuilder for the given table.
//
// Example:
//
//	q := query.Delete("users").
//		Where("id = ?", 42)
//
//	sql, args := q.Build()
//	 MySQL:    "DELETE FROM users WHERE id = ?"
//	 Postgres: "DELETE FROM users WHERE id = $1"
//	 args: [42]
func Delete(table string) *DeleteBuilder {
	return &DeleteBuilder{
		table: table,
	}
}

// Where adds a WHERE clause to the DELETE query.
// Multiple calls are joined with AND.
func (b *DeleteBuilder) Where(conds ...Condition) *DeleteBuilder {
	b.where = append(b.where, conds...)
	return b
}

// Build assembles the SQL DELETE query string and returns it with args.
// It rewrites placeholders depending on the active dialect.
func (b *DeleteBuilder) Build() (string, []any) {
	var sql strings.Builder
	sql.WriteString("DELETE FROM ")
	sql.WriteString(b.table)

	if len(b.where) > 0 {
		sql.WriteString(" WHERE ")
		parts := make([]string, len(b.where))
		for i, cond := range b.where {
			expr := cond.Expr
			for j := 0; j < len(cond.Args); j++ {
				b.placeholderIndex++
				expr = strings.Replace(expr, "?", GetDialect().Placeholder(b.placeholderIndex), 1)
			}
			parts[i] = expr
			b.args = append(b.args, cond.Args...)
		}
		sql.WriteString(strings.Join(parts, " AND "))
	}

	return sql.String(), b.args
}
