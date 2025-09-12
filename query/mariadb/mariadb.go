// Package mariadb registers the PostgreSQL dialect for the OCA query builder.
//
// Importing this package will automatically set the query dialect to
// PostgreSQL via the init function:
//
//	import _ "github.com/mhdiiilham/oca/query/mariadb"
//
// After importing, all queries built using the query package will be
// generated using PostgreSQL syntax without requiring manual dialect setup.
package mariadb

import "github.com/mhdiiilham/oca/query"

// init registers the PostgreSQL dialect as the default dialect
// for the query builder when this package is imported.
func init() {
	query.SetDialect(query.MariaDBDialect{})
}
