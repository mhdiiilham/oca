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

func TestRepository_Delete_WithWhere(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[Todo](db)

	mock.ExpectExec(`DELETE FROM todos WHERE id = \?`).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(context.Background(), oca.Where(query.C("id").Eq(int64(1))))
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete_WithMultipleWhere(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[Todo](db)

	mock.ExpectExec(`DELETE FROM todos WHERE id = \? AND title = \?`).
		WithArgs(int64(5), "Task 5").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(context.Background(),
		oca.Where(query.C("id").Eq(int64(5))),
		oca.Where(query.C("title").Eq("Task 5")),
	)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete_ExecError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[Todo](db)

	mock.ExpectExec(`DELETE FROM todos WHERE id = \?`).
		WithArgs(int64(42)).
		WillReturnError(sql.ErrConnDone)

	err := repo.Delete(context.Background(), oca.Where(query.C("id").Eq(int64(42))))
	assert.ErrorIs(t, err, sql.ErrConnDone)
	assert.NoError(t, mock.ExpectationsWereMet())
}
