package query

// joinClause stores info about a JOIN clause
type joinClause struct {
	kind  string
	table string
	on    string
}

// Join adds an INNER JOIN clause.
func (b *Builder) Join(table, on string) *Builder {
	b.joins = append(b.joins, joinClause{"INNER JOIN", table, on})
	return b
}

// LeftJoin adds a LEFT JOIN clause.
func (b *Builder) LeftJoin(table, on string) *Builder {
	b.joins = append(b.joins, joinClause{"LEFT JOIN", table, on})
	return b
}

// RightJoin adds a RIGHT JOIN clause.
func (b *Builder) RightJoin(table, on string) *Builder {
	b.joins = append(b.joins, joinClause{"RIGHT JOIN", table, on})
	return b
}

// FullJoin adds a FULL JOIN clause.
func (b *Builder) FullJoin(table, on string) *Builder {
	b.joins = append(b.joins, joinClause{"FULL JOIN", table, on})
	return b
}
