package query

import (
	"fmt"
	"strings"
)

// RawSQL wraps a SQL literal so it can be inserted directly into the query
// without being parameterized.
type RawSQL string

// Raw creates a RawSQL instance.
// Example: query.Raw("NOW()")
func Raw(val string) RawSQL {
	return RawSQL(val)
}

// InsertBuilder builds SQL INSERT statements.
type InsertBuilder struct {
	table     string
	columns   []string
	values    [][]any
	returning []string
	dialect   Dialect // nil → falls back to global
}

// InsertInto creates a new InsertBuilder for the given table.
func InsertInto(table string) *InsertBuilder {
	return &InsertBuilder{table: table}
}

// WithDialect overrides the SQL dialect for this builder only.
func (b *InsertBuilder) WithDialect(d Dialect) *InsertBuilder {
	b.dialect = d
	return b
}

func (b *InsertBuilder) getDialect() Dialect {
	if b.dialect != nil {
		return b.dialect
	}
	return GetDialect()
}

// Columns sets the columns for the INSERT statement.
func (b *InsertBuilder) Columns(cols ...string) *InsertBuilder {
	b.columns = append(b.columns, cols...)
	return b
}

// Values adds a row of values. Supports RawSQL for literals like NOW().
// Each call adds one row; call multiple times for multi-row inserts.
func (b *InsertBuilder) Values(vals ...any) *InsertBuilder {
	row := make([]any, len(vals))
	copy(row, vals)
	b.values = append(b.values, row)
	return b
}

// Returning specifies columns to return after insert (PostgreSQL).
func (b *InsertBuilder) Returning(cols ...string) *InsertBuilder {
	b.returning = append(b.returning, cols...)
	return b
}

// ReturningID is shorthand for Returning("id").
func (b *InsertBuilder) ReturningID() *InsertBuilder {
	return b.Returning("id")
}

// ToSQL builds the final INSERT query and returns the SQL string and arguments.
func (b *InsertBuilder) ToSQL() (string, []any) {
	if b.table == "" || len(b.columns) == 0 || len(b.values) == 0 {
		return "", nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES ",
		b.table, strings.Join(b.columns, ", ")))

	var args []any
	var placeholders []string
	placeholderIdx := 0

	for _, row := range b.values {
		rowPlaceholders := make([]string, len(row))
		for i, v := range row {
			switch v := v.(type) {
			case RawSQL:
				rowPlaceholders[i] = string(v)
			default:
				placeholderIdx++
				rowPlaceholders[i] = b.getDialect().Placeholder(placeholderIdx)
				args = append(args, v)
			}
		}
		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ", ")))
	}

	sb.WriteString(strings.Join(placeholders, ", "))

	if len(b.returning) > 0 {
		sb.WriteString(" RETURNING ")
		sb.WriteString(strings.Join(b.returning, ", "))
	}

	return sb.String(), args
}
