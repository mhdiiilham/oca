package oca

import "context"

// GenericStore defines the full CRUD interface for a repository backed by model T.
// T can be any struct annotated with `db` tags for column mapping.
type GenericStore[T any] interface {
	// Insert inserts one record and populates auto-increment / schema-default fields.
	//
	//	u := &User{Name: "Alice"}
	//	err := repo.Insert(ctx, u) // u.ID set after insert
	Insert(ctx context.Context, entity *T) error

	// InsertMany inserts multiple records in a single query.
	// Auto-increment fields are not populated back; use Insert for that.
	//
	//	err := repo.InsertMany(ctx, []User{{Name: "Alice"}, {Name: "Bob"}})
	InsertMany(ctx context.Context, entities []T) error

	// Finds retrieves all matching records.
	//
	//	todos, err := repo.Finds(ctx, oca.Where(query.C("done").Eq(false)))
	Finds(ctx context.Context, filter ...FilterOptions) ([]T, error)

	// FindOne retrieves one record. Returns sql.ErrNoRows when nothing matches.
	//
	//	todo, err := repo.FindOne(ctx, oca.Where(query.C("id").Eq(28)))
	FindOne(ctx context.Context, opts ...FilterOptions) (*T, error)

	// Update writes all non-auto, non-pk fields back to the database.
	// Pointer fields that are nil are skipped (partial update).
	// Defaults to WHERE pk = entity.PK; pass Where() opts to override.
	//
	//	u.Name = "Bob"
	//	err := repo.Update(ctx, u) // UPDATE users SET name=? WHERE id=?
	Update(ctx context.Context, entity *T, opts ...FilterOptions) error

	// Delete removes rows matching the filter. At least one Where() condition
	// is required to prevent accidental full-table deletion.
	// Use DeleteAll to intentionally remove every row.
	//
	//	err := repo.Delete(ctx, oca.Where(query.C("id").Eq(u.ID)))
	Delete(ctx context.Context, opts ...FilterOptions) error

	// DeleteAll removes every row in the table. Use with care.
	DeleteAll(ctx context.Context) error

	// Tx runs fn inside a single database transaction, committing on success
	// and rolling back on any error returned by fn.
	//
	//	err := repo.Tx(ctx, func(tx oca.GenericStore[User]) error {
	//	    if err := tx.Insert(ctx, &u); err != nil { return err }
	//	    return tx.Update(ctx, &other)
	//	})
	Tx(ctx context.Context, fn func(store GenericStore[T]) error) error
}
