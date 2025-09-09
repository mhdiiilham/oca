package oca_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mhdiiilham/oca"
	"github.com/stretchr/testify/assert"
)

// SimpleUser has no auto-increment field
type SimpleUser struct {
	Name string `db:"name"`
}

// AutoUser has an auto-increment primary key
type AutoUser struct {
	ID   int64  `db:"id,pk,auto"`
	Name string `db:"name"`
}

// TodoWithDefault demonstrates schema default:now()
type TodoWithDefault struct {
	ID        int64     `db:"id,pk,auto"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at" schema:"default:now()"`
}

func TestRepository_Insert_NoAuto(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := oca.NewRepository[SimpleUser](db)

	mock.ExpectExec(`INSERT INTO simpleusers \(name\) VALUES \(\?\)`).
		WithArgs("Alice").
		WillReturnResult(sqlmock.NewResult(0, 1))

	u := &SimpleUser{Name: "Alice"}
	err = repo.Insert(context.Background(), u)

	assert.NoError(t, err)
	assert.Equal(t, "Alice", u.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Insert_WithAuto(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := oca.NewRepository[AutoUser](db)

	mock.ExpectQuery(`INSERT INTO autousers \(name\) VALUES \(\?\) RETURNING id`).
		WithArgs("Bob").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

	u := &AutoUser{Name: "Bob"}
	err = repo.Insert(context.Background(), u)

	assert.NoError(t, err)
	assert.Equal(t, int64(42), u.ID)
	assert.Equal(t, "Bob", u.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Insert_ErrorOnExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := oca.NewRepository[SimpleUser](db)

	mock.ExpectExec(`INSERT INTO simpleusers \(name\) VALUES \(\?\)`).
		WithArgs("Charlie").
		WillReturnError(sql.ErrConnDone)

	u := &SimpleUser{Name: "Charlie"}
	err = repo.Insert(context.Background(), u)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Insert_ErrorOnReturning(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := oca.NewRepository[AutoUser](db)

	mock.ExpectQuery(`INSERT INTO autousers \(name\) VALUES \(\?\) RETURNING id`).
		WithArgs("Dave").
		WillReturnError(sql.ErrNoRows)

	u := &AutoUser{Name: "Dave"}
	err = repo.Insert(context.Background(), u)

	assert.Error(t, err)
	assert.Equal(t, int64(0), u.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Insert_DefaultNow(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := oca.NewRepository[TodoWithDefault](db)

	// Expect INSERT with NOW() and RETURNING id
	mock.ExpectQuery(`INSERT INTO todowithdefaults \(title, created_at\) VALUES \(\?, NOW\(\)\) RETURNING id`).
		WithArgs("Task 1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	todo := &TodoWithDefault{Title: "Task 1"}
	err = repo.Insert(context.Background(), todo)

	assert.NoError(t, err)
	assert.Equal(t, "Task 1", todo.Title)
	assert.Equal(t, int64(1), todo.ID) // make sure ID was set
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Insert_MultipleAutoFields(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	type Todo struct {
		ID        int64     `db:"id,pk,auto"`
		Title     string    `db:"title"`
		CreatedAt time.Time `db:"created_at,auto"`
	}

	repo := oca.NewRepository[Todo](db)

	// Simulate DB returning id and created_at
	now := time.Now()
	mock.ExpectQuery(`INSERT INTO todos \(title\) VALUES \(\?\) RETURNING id, created_at`).
		WithArgs("Task Multi").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(100, now))

	todo := &Todo{Title: "Task Multi"}
	err = repo.Insert(context.Background(), todo)

	assert.NoError(t, err)
	assert.Equal(t, int64(100), todo.ID, "ID should be returned by DB")
	assert.Equal(t, now, todo.CreatedAt, "CreatedAt should be returned by DB")
	assert.Equal(t, "Task Multi", todo.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}
