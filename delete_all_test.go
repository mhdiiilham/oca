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

func TestRepository_Delete_WithoutWhere_ReturnsError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[Todo](db)

	err := repo.Delete(context.Background())
	assert.ErrorContains(t, err, "oca: Delete requires")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[Todo](db)

	mock.ExpectExec(`DELETE FROM todos`).
		WillReturnResult(sqlmock.NewResult(0, 5))

	err := repo.DeleteAll(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteAll_ExecError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[Todo](db)

	mock.ExpectExec(`DELETE FROM todos`).
		WillReturnError(sql.ErrConnDone)

	err := repo.DeleteAll(context.Background())
	assert.ErrorIs(t, err, sql.ErrConnDone)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_WithDialect_Postgres(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[AutoUser](db, oca.WithDialect[AutoUser](query.PostgresDialect{}))

	mock.ExpectQuery(`INSERT INTO autousers \(name\) VALUES \(\$1\) RETURNING id`).
		WithArgs("Alice").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	u := &AutoUser{Name: "Alice"}
	err := repo.Insert(context.Background(), u)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), u.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_WithDialect_Postgres_Finds(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[AutoUser](db, oca.WithDialect[AutoUser](query.PostgresDialect{}))

	mock.ExpectQuery(`SELECT id, name FROM autousers WHERE id = \$1 LIMIT \$2`).
		WithArgs(int64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(int64(1), "Alice"))

	u, err := repo.FindOne(context.Background(), oca.Where(query.C("id").Eq(int64(1))))
	assert.NoError(t, err)
	assert.Equal(t, int64(1), u.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
