package query

import (
	"strings"
)

// DeleteBuilder builds SQL DELETE queries with dialect support.
type DeleteBuilder struct {
	table            string
	where            []Condition
	args             []any
	dialect          Dialect // nil → falls back to global
	placeholderIndex int
}

// Delete creates a new DeleteBuilder for the given table.
//
// Example:
//
//	q := query.Delete("users").Where(query.C("id").Eq(42))
//	sql, args := q.Build()
//	 MySQL:    "DELETE FROM users WHERE id = ?"
//	 Postgres: "DELETE FROM users WHERE id = $1"
func Delete(table string) *DeleteBuilder {
	return &DeleteBuilder{table: table}
}

// WithDialect overrides the SQL dialect for this builder only.
func (b *DeleteBuilder) WithDialect(d Dialect) *DeleteBuilder {
	b.dialect = d
	return b
}

func (b *DeleteBuilder) getDialect() Dialect {
	if b.dialect != nil {
		return b.dialect
	}
	return GetDialect()
}

// Where adds a WHERE clause. Multiple calls are joined with AND.
func (b *DeleteBuilder) Where(conds ...Condition) *DeleteBuilder {
	b.where = append(b.where, conds...)
	return b
}

// Build assembles the SQL DELETE query string and returns it with args.
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
				expr = strings.Replace(expr, "?", b.getDialect().Placeholder(b.placeholderIndex), 1)
			}
			parts[i] = expr
			b.args = append(b.args, cond.Args...)
		}
		sql.WriteString(strings.Join(parts, " AND "))
	}

	return sql.String(), b.args
}
