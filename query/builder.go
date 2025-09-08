package query

import (
	"strings"
)

// Builder builds SQL queries in a fluent style.
type Builder struct {
	table   string
	columns []string
	where   []string
	args    []any
	order   string
	limit   int
	offset  int
}

// From starts a new query builder for the given table.
func From(table string) *Builder {
	return &Builder{table: table, limit: -1, offset: -1}
}

// Select adds columns to the SELECT clause.
// If no columns are added, "*" will be used.
func (b *Builder) Select(cols ...string) *Builder {
	b.columns = append(b.columns, cols...)
	return b
}

// Where adds a WHERE condition with optional arguments.
// Multiple calls to Where will be combined with AND.
func (b *Builder) Where(cond string, args ...any) *Builder {
	b.where = append(b.where, cond)
	b.args = append(b.args, args...)
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

// Build assembles the SQL query string and returns it
// along with the bound arguments.
func (b *Builder) Build() (string, []any) {
	cols := "*"
	if len(b.columns) > 0 {
		cols = strings.Join(b.columns, ", ")
	}

	sql := strings.Builder{}
	sql.WriteString("SELECT ")
	sql.WriteString(cols)
	sql.WriteString(" FROM ")
	sql.WriteString(b.table)

	if len(b.where) > 0 {
		sql.WriteString(" WHERE ")
		sql.WriteString(strings.Join(b.where, " AND "))
	}

	if b.order != "" {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(b.order)
	}

	if b.limit >= 0 {
		sql.WriteString(" LIMIT ?")
		b.args = append(b.args, b.limit)
	}

	if b.offset >= 0 {
		sql.WriteString(" OFFSET ?")
		b.args = append(b.args, b.offset)
	}

	return sql.String(), b.args
}
