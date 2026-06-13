package oca_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mhdiiilham/oca"
	"github.com/stretchr/testify/assert"
)

// PartialUser has one pointer field to test partial-update behaviour.
type PartialUser struct {
	ID    int64   `db:"id,pk,auto"`
	Name  *string `db:"name"`
	Email string  `db:"email"`
}

func ptr(s string) *string { return &s }

func TestRepository_Update_NilPointerFieldSkipped(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[PartialUser](db)

	// Name is nil → only email appears in SET
	mock.ExpectExec(`UPDATE partialusers SET email = \? WHERE id = \?`).
		WithArgs("alice@example.com", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	u := &PartialUser{ID: 1, Name: nil, Email: "alice@example.com"}
	err := repo.Update(context.Background(), u)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_NonNilPointerIncluded(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[PartialUser](db)

	// Name is non-nil → both name and email appear in SET
	mock.ExpectExec(`UPDATE partialusers SET name = \?, email = \? WHERE id = \?`).
		WithArgs("Alice", "alice@example.com", int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	u := &PartialUser{ID: 2, Name: ptr("Alice"), Email: "alice@example.com"}
	err := repo.Update(context.Background(), u)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
