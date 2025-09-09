package oca

import "context"

// GenericStore defines a generic interface for basic CRUD operations on any model T.
// This allows repositories to be abstracted and easily mocked for testing.
//
// T can be any struct that represents a database model.
type GenericStore[T any] interface {
	// Insert inserts a new record into the database and sets auto-increment ID or other auto fields if applicable.
	// It also respects schema defaults, e.g., `schema:"default:now()"`.
	//
	// Example:
	//
	//	type User struct {
	//	    ID        int64  `db:"id,pk,auto"`
	//	    Name      string `db:"name"`
	//	    CreatedAt time.Time `db:"created_at" schema:"default:now()"`
	//	}
	//
	//	repo := NewRepository[User](db)
	//	u := &User{Name: "Alice"}
	//	err := repo.Insert(ctx, u)
	//
	// u.ID and u.CreatedAt will be automatically set if applicable.
	Insert(ctx context.Context, entity *T) error
}
