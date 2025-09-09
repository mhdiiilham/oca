package query

import (
	"fmt"
	"sync"
)

// Dialect defines the behavior that differs across SQL databases.
// For example, how parameter placeholders are represented.
type Dialect interface {
	// Placeholder returns the placeholder string for the given index.
	// Example: $1 for PostgreSQL, ? for MySQL/MariaDB.
	Placeholder(index int) string
	// Name returns the name of the dialect (for debugging/logging).
	Name() string
}

// ------------------
// Dialect Implementations
// ------------------

// MySQLDialect uses "?" placeholders.
type MySQLDialect struct{}

// Placeholder returns "?" for all indexes.
func (d MySQLDialect) Placeholder(_ int) string {
	return "?"
}

// Name returns the dialect name.
func (d MySQLDialect) Name() string {
	return "mysql"
}

// MariaDBDialect behaves the same as MySQL for placeholders.
type MariaDBDialect struct{}

// Placeholder returns "?" for all indexes.
func (d MariaDBDialect) Placeholder(_ int) string {
	return "?"
}

// Name returns the dialect name.
func (d MariaDBDialect) Name() string {
	return "mariadb"
}

// PostgresDialect uses "$1, $2, ..." placeholders.
type PostgresDialect struct{}

// Placeholder returns "$n".
func (d PostgresDialect) Placeholder(i int) string {
	return fmt.Sprintf("$%d", i)
}

// Name returns the dialect name.
func (d PostgresDialect) Name() string {
	return "postgresql"
}

// ------------------
// Global Dialect Management
// ------------------

var (
	mu            sync.RWMutex
	globalDialect Dialect = MySQLDialect{} // default is MySQL
)

// SetDialect sets the global SQL dialect for all queries.
// Example: query.SetDialect(query.PostgresDialect{})
func SetDialect(d Dialect) {
	mu.Lock()
	defer mu.Unlock()
	globalDialect = d
}

// GetDialect returns the current global dialect.
func GetDialect() Dialect {
	mu.RLock()
	defer mu.RUnlock()
	return globalDialect
}
