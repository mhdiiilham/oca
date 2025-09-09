package query

import (
	"strings"
)

// Builder builds SQL SELECT queries in a fluent DSL style.
// It supports SELECT, WHERE, ORDER BY, LIMIT, OFFSET.
type Builder struct {
	table            string
	columns          []string
	where            []Condition
	args             []any
	order            string
	limit            int
	offset           int
	joins            []joinClause
	placeholderIndex int // tracks placeholders for dialects
}

// From creates a new Builder for a given table.
func From(table string) *Builder {
	return &Builder{table: table, limit: -1, offset: -1}
}

// Select specifies the columns to select.
func (b *Builder) Select(cols ...string) *Builder {
	b.columns = cols
	return b
}

// Where adds one or more conditions. Multiple calls are combined with AND.
func (b *Builder) Where(conds ...Condition) *Builder {
	b.where = append(b.where, conds...)
	return b
}

// OrderBy sets the ORDER BY clause.
func (b *Builder) OrderBy(order string) *Builder {
	b.order = order
	return b
}

// Limit sets the LIMIT clause.
func (b *Builder) Limit(limit int) *Builder {
	b.limit = limit
	return b
}

// Offset sets the OFFSET clause.
func (b *Builder) Offset(offset int) *Builder {
	b.offset = offset
	return b
}

// Build assembles the SQL string and returns args.
func (b *Builder) Build() (string, []any) {
	b.args = nil
	b.placeholderIndex = 0

	cols := "*"
	if len(b.columns) > 0 {
		cols = strings.Join(b.columns, ", ")
	}

	sql := strings.Builder{}
	sql.WriteString("SELECT ")
	sql.WriteString(cols)
	sql.WriteString(" FROM ")
	sql.WriteString(b.table)

	for _, j := range b.joins {
		sql.WriteString(" ")
		sql.WriteString(j.kind)
		sql.WriteString(" ")
		sql.WriteString(j.table)
		sql.WriteString(" ON ")
		sql.WriteString(j.on)
	}

	// WHERE
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

	// ORDER BY
	if b.order != "" {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(b.order)
	}

	// LIMIT
	if b.limit >= 0 {
		b.placeholderIndex++
		sql.WriteString(" LIMIT ")
		sql.WriteString(GetDialect().Placeholder(b.placeholderIndex))
		b.args = append(b.args, b.limit)
	}

	// OFFSET
	if b.offset >= 0 {
		b.placeholderIndex++
		sql.WriteString(" OFFSET ")
		sql.WriteString(GetDialect().Placeholder(b.placeholderIndex))
		b.args = append(b.args, b.offset)
	}

	return sql.String(), b.args
}
