package query

// Select adds columns to the SELECT clause.
// If no columns are added, "*" will be used.
func (b *Builder) Select(cols ...string) *Builder {
	b.columns = append(b.columns, cols...)
	return b
}
