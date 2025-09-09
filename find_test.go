package oca_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mhdiiilham/oca"
	"github.com/mhdiiilham/oca/query"
	"github.com/stretchr/testify/assert"
)

type Todo struct {
	ID        int64     `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
}

func TestRepository_Find(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := oca.NewRepository[Todo](db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "title", "created_at"}).
		AddRow(1, "Task 1", now).
		AddRow(2, "Task 2", now.Add(time.Hour))

	mock.ExpectQuery(`SELECT id, title, created_at FROM todos WHERE title = \? ORDER BY created_at DESC LIMIT \?`).
		WithArgs("Task 1", 10).
		WillReturnRows(rows)

	todos, err := repo.Find(context.Background(),
		oca.OrderBy("created_at DESC"),
		oca.Limit(10),
		oca.Offset(0),
		oca.Where(
			query.C("title").Eq("Task 1"),
		),
	)

	assert.NoError(t, err)
	assert.Len(t, todos, 2)
	assert.Equal(t, int64(1), todos[0].ID)
	assert.Equal(t, "Task 1", todos[0].Title)
	assert.Equal(t, int64(2), todos[1].ID)
	assert.Equal(t, "Task 2", todos[1].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_FindOne_ByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := oca.NewRepository[Todo](db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "title", "created_at"}).
		AddRow(1, "Task 1", now)

	// Expect SELECT with WHERE id = ?
	mock.ExpectQuery(`SELECT id, title, created_at FROM todos WHERE id = \? LIMIT \?`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	todo, err := repo.FindOne(
		context.Background(),
		oca.Where(query.C("id").Eq(1)),
	)
	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, int64(1), todo.ID)
	assert.Equal(t, "Task 1", todo.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_FindOne_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := oca.NewRepository[Todo](db)

	// Expect SELECT but return no rows
	rows := sqlmock.NewRows([]string{"id", "title", "created_at"})

	mock.ExpectQuery(`SELECT id, title, created_at FROM todos WHERE id = \? LIMIT \?`).
		WithArgs(999, 1).
		WillReturnRows(rows)

	todo, err := repo.FindOne(
		context.Background(),
		oca.Where(query.C("id").Eq(999)),
	)

	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Nil(t, todo)
	assert.NoError(t, mock.ExpectationsWereMet())
}
