package oca

import (
	"database/sql"
)

// Repository provides basic CRUD operations for any model T.
// If T implements Tabler, its TableName() will be used.
// Otherwise, the struct name is used as the table name (lowercased + "s").
type Repository[T any] struct {
	db *sql.DB
}

// NewRepository returns a new generic repository.
func NewRepository[T any](db *sql.DB) GenericStore[T] {
	return &Repository[T]{db: db}
}
