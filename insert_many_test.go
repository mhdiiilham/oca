package oca_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mhdiiilham/oca"
	"github.com/stretchr/testify/assert"
)

func TestRepository_InsertMany_NoAuto(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[SimpleUser](db)

	mock.ExpectExec(`INSERT INTO simpleusers \(name\) VALUES \(\?\), \(\?\)`).
		WithArgs("Alice", "Bob").
		WillReturnResult(sqlmock.NewResult(0, 2))

	err := repo.InsertMany(context.Background(), []SimpleUser{{Name: "Alice"}, {Name: "Bob"}})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_InsertMany_ThreeRows(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[SimpleUser](db)

	mock.ExpectExec(`INSERT INTO simpleusers \(name\) VALUES \(\?\), \(\?\), \(\?\)`).
		WithArgs("Alice", "Bob", "Carol").
		WillReturnResult(sqlmock.NewResult(0, 3))

	err := repo.InsertMany(context.Background(), []SimpleUser{
		{Name: "Alice"},
		{Name: "Bob"},
		{Name: "Carol"},
	})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_InsertMany_Empty(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[SimpleUser](db)

	err := repo.InsertMany(context.Background(), []SimpleUser{})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_InsertMany_ExecError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[SimpleUser](db)

	mock.ExpectExec(`INSERT INTO simpleusers \(name\) VALUES \(\?\)`).
		WithArgs("Alice").
		WillReturnError(sql.ErrConnDone)

	err := repo.InsertMany(context.Background(), []SimpleUser{{Name: "Alice"}})
	assert.ErrorIs(t, err, sql.ErrConnDone)
	assert.NoError(t, mock.ExpectationsWereMet())
}
