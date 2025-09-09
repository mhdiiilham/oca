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

type SimpleUser struct {
	Name string `db:"name"`
}

type AutoUser struct {
	ID   int64  `db:"id,pk,auto"`
	Name string `db:"name"`
}

type TodoWithDefault struct {
	ID        int64     `db:"id,pk,auto"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at" schema:"default:now()"`
}

func TestRepository_Insert_NoAuto(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[SimpleUser](db)

	mock.ExpectExec(`INSERT INTO simpleusers \(name\) VALUES \(\?\)`).
		WithArgs("Alice").
		WillReturnResult(sqlmock.NewResult(0, 1))

	u := &SimpleUser{Name: "Alice"}
	err := repo.Insert(context.Background(), u)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", u.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Insert_WithAuto(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[AutoUser](db)

	mock.ExpectQuery(`INSERT INTO autousers \(name\) VALUES \(\?\) RETURNING id`).
		WithArgs("Bob").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

	u := &AutoUser{Name: "Bob"}
	err := repo.Insert(context.Background(), u)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), u.ID)
	assert.Equal(t, "Bob", u.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Insert_ErrorOnExec(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[SimpleUser](db)

	mock.ExpectExec(`INSERT INTO simpleusers \(name\) VALUES \(\?\)`).
		WithArgs("Charlie").
		WillReturnError(sql.ErrConnDone)

	u := &SimpleUser{Name: "Charlie"}
	err := repo.Insert(context.Background(), u)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Insert_ErrorOnReturning(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[AutoUser](db)

	mock.ExpectQuery(`INSERT INTO autousers \(name\) VALUES \(\?\) RETURNING id`).
		WithArgs("Dave").
		WillReturnError(sql.ErrNoRows)

	u := &AutoUser{Name: "Dave"}
	err := repo.Insert(context.Background(), u)
	assert.Error(t, err)
	assert.Equal(t, int64(0), u.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Insert_DefaultNow(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[TodoWithDefault](db)

	mock.ExpectQuery(`INSERT INTO todowithdefaults \(title, created_at\) VALUES \(\?, NOW\(\)\) RETURNING id`).
		WithArgs("Task 1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	todo := &TodoWithDefault{Title: "Task 1"}
	err := repo.Insert(context.Background(), todo)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), todo.ID)
	assert.Equal(t, "Task 1", todo.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Insert_MultipleAutoFields(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	type Todo struct {
		ID        int64     `db:"id,pk,auto"`
		Title     string    `db:"title"`
		CreatedAt time.Time `db:"created_at,auto"`
	}

	repo := oca.NewRepository[Todo](db)
	now := time.Now()

	mock.ExpectQuery(`INSERT INTO todos \(title\) VALUES \(\?\) RETURNING id, created_at`).
		WithArgs("Task Multi").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(100, now))

	todo := &Todo{Title: "Task Multi"}
	err := repo.Insert(context.Background(), todo)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), todo.ID)
	assert.Equal(t, now, todo.CreatedAt)
	assert.Equal(t, "Task Multi", todo.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Insert_DefaultNow_WithAutoField(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	type Todo struct {
		ID        int64     `db:"id,pk,auto"`
		Title     string    `db:"title"`
		CreatedAt time.Time `db:"created_at" schema:"default:now()"`
	}

	repo := oca.NewRepository[Todo](db)

	mock.ExpectQuery(`INSERT INTO todos \(title, created_at\) VALUES \(\?, NOW\(\)\) RETURNING id`).
		WithArgs("Task Default").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(200))

	todo := &Todo{Title: "Task Default"}
	err := repo.Insert(context.Background(), todo)
	assert.NoError(t, err)
	assert.Equal(t, int64(200), todo.ID)
	assert.Equal(t, "Task Default", todo.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}
