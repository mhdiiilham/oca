package oca

// Tabler is implemented by any struct that wants to define
// its own database mapping. It provides the table name used in queries.
//
// Example:
//
//	func (User) TableName() string   { return "users" }
type Tabler interface {
	// TableName returns the name of the database table.
	TableName() string
}
