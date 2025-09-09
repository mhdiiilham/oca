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
	values    [][]interface{}
	returning []string
}

// InsertInto creates a new InsertBuilder for the given table.
// Example: query.InsertInto("users")
func InsertInto(table string) *InsertBuilder {
	return &InsertBuilder{table: table}
}

// Columns sets the columns for the INSERT statement.
// Example: .Columns("id", "name")
func (b *InsertBuilder) Columns(cols ...string) *InsertBuilder {
	b.columns = append(b.columns, cols...)
	return b
}

// Values adds a row of values to insert.
// Supports RawSQL for literals.
// Example: .Values(1, "John", query.Raw("NOW()"))
func (b *InsertBuilder) Values(vals ...interface{}) *InsertBuilder {
	row := make([]interface{}, len(vals))
	copy(row, vals) // <-- simpler and linter-friendly
	b.values = append(b.values, row)
	return b
}

// Returning specifies columns to return (Postgres style).
// Example: .Returning("id", "created_at")
func (b *InsertBuilder) Returning(cols ...string) *InsertBuilder {
	b.returning = append(b.returning, cols...)
	return b
}

// ReturningID is a convenience method for "RETURNING id".
func (b *InsertBuilder) ReturningID() *InsertBuilder {
	return b.Returning("id")
}

// ToSQL builds the final INSERT query and returns the SQL string and arguments.
// Example output: "INSERT INTO users (id, name) VALUES (?, ?) RETURNING id", [1, "John"]
func (b *InsertBuilder) ToSQL() (string, []interface{}) {
	if b.table == "" || len(b.columns) == 0 || len(b.values) == 0 {
		return "", nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES ",
		b.table, strings.Join(b.columns, ", ")))

	var args []interface{}
	var placeholders []string

	for _, row := range b.values {
		rowPlaceholders := make([]string, len(row))
		for i, v := range row {
			switch v := v.(type) {
			case RawSQL:
				rowPlaceholders[i] = string(v) // literal SQL
			default:
				rowPlaceholders[i] = "?"
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
