package oca

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type userTest struct {
	ID        int64     `db:"id,pk,auto"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func Test_scanRows(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "created_at"}).
		AddRow(1, "Alice", time.Now()).
		AddRow(2, "Bob", time.Now())

	mock.ExpectQuery("SELECT id, name, created_at FROM users").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT id, name, created_at FROM users")
	assert.NoError(t, err)

	results, err := scanRows[userTest](sqlRows)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, int64(1), results[0].ID)
	assert.Equal(t, "Alice", results[0].Name)
	assert.Equal(t, int64(2), results[1].ID)
	assert.Equal(t, "Bob", results[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_scanRowSingle(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery("SELECT id, name, created_at FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).AddRow(1, "Alice", now))

	row := db.QueryRow("SELECT id, name, created_at FROM users WHERE id = ?", 1)

	columns := []string{"id", "name", "created_at"}
	colMap := map[string]int{"id": 0, "name": 1, "created_at": 2}

	result, err := scanRowSingle[userTest](row, columns, colMap)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Alice", result.Name)
	assert.Equal(t, now, result.CreatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_scanAutoFields(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	autoFields := []FieldMeta{
		{Column: "id", Index: 0},
	}

	mock.ExpectQuery("INSERT INTO users \\(name\\) VALUES \\(\\?\\) RETURNING id").
		WithArgs("Alice").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	user := &userTest{Name: "Alice"}

	err := scanAutoFields(context.Background(), db, "INSERT INTO users (name) VALUES (?) RETURNING id", []interface{}{"Alice"}, user, autoFields)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
