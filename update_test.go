package oca_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mhdiiilham/oca"
	"github.com/mhdiiilham/oca/query"
	"github.com/stretchr/testify/assert"
)

// UpdateUser is a model with pk+auto ID and two mutable fields.
type UpdateUser struct {
	ID    int64  `db:"id,pk,auto"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func TestRepository_Update_ByPrimaryKey(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[UpdateUser](db)

	mock.ExpectExec(`UPDATE updateusers SET name = \?, email = \? WHERE id = \?`).
		WithArgs("Alice", "alice@example.com", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	u := &UpdateUser{ID: 1, Name: "Alice", Email: "alice@example.com"}
	err := repo.Update(context.Background(), u)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_WithWhereOpt(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[UpdateUser](db)

	mock.ExpectExec(`UPDATE updateusers SET name = \?, email = \? WHERE name = \?`).
		WithArgs("Bob", "bob@example.com", "Alice").
		WillReturnResult(sqlmock.NewResult(0, 2))

	u := &UpdateUser{Name: "Bob", Email: "bob@example.com"}
	err := repo.Update(context.Background(), u, oca.Where(query.C("name").Eq("Alice")))
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_ExecError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[UpdateUser](db)

	mock.ExpectExec(`UPDATE updateusers SET name = \?, email = \? WHERE id = \?`).
		WithArgs("Charlie", "", int64(99)).
		WillReturnError(sql.ErrConnDone)

	u := &UpdateUser{ID: 99, Name: "Charlie"}
	err := repo.Update(context.Background(), u)
	assert.ErrorIs(t, err, sql.ErrConnDone)
	assert.NoError(t, mock.ExpectationsWereMet())
}
