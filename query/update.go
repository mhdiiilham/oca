package query

import (
	"strings"
)

type setClause struct {
	col string
	val any
}

// UpdateBuilder builds SQL UPDATE queries with dialect support.
type UpdateBuilder struct {
	table            string
	sets             []setClause
	where            []Condition
	args             []any
	dialect          Dialect // nil → falls back to global
	placeholderIndex int
}

// Update creates a new UpdateBuilder for the given table.
//
// Example:
//
//	q := query.Update("users").
//		Set("name", "Alice").
//		Where(query.C("id").Eq(1))
//
//	sql, args := q.Build()
//	 MySQL:    "UPDATE users SET name = ? WHERE id = ?"
//	 Postgres: "UPDATE users SET name = $1 WHERE id = $2"
func Update(table string) *UpdateBuilder {
	return &UpdateBuilder{table: table}
}

// WithDialect overrides the SQL dialect for this builder only.
func (b *UpdateBuilder) WithDialect(d Dialect) *UpdateBuilder {
	b.dialect = d
	return b
}

func (b *UpdateBuilder) getDialect() Dialect {
	if b.dialect != nil {
		return b.dialect
	}
	return GetDialect()
}

// Set adds a column = value assignment to the SET clause.
func (b *UpdateBuilder) Set(col string, val any) *UpdateBuilder {
	b.sets = append(b.sets, setClause{col: col, val: val})
	return b
}

// Where adds WHERE conditions. Multiple calls are joined with AND.
func (b *UpdateBuilder) Where(conds ...Condition) *UpdateBuilder {
	b.where = append(b.where, conds...)
	return b
}

// Build assembles the SQL UPDATE query string and returns it with args.
// Placeholder numbering is shared across SET and WHERE so Postgres $N is correct.
func (b *UpdateBuilder) Build() (string, []any) {
	var sql strings.Builder
	sql.WriteString("UPDATE ")
	sql.WriteString(b.table)
	sql.WriteString(" SET ")

	setParts := make([]string, len(b.sets))
	for i, s := range b.sets {
		b.placeholderIndex++
		setParts[i] = s.col + " = " + b.getDialect().Placeholder(b.placeholderIndex)
		b.args = append(b.args, s.val)
	}
	sql.WriteString(strings.Join(setParts, ", "))

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
