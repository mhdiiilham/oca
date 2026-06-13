package oca_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mhdiiilham/oca"
	"github.com/mhdiiilham/oca/query"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Tx_Commit(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[AutoUser](db)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO autousers \(name\) VALUES \(\?\) RETURNING id`).
		WithArgs("Alice").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))
	mock.ExpectCommit()

	var inserted AutoUser
	err := repo.Tx(context.Background(), func(tx oca.GenericStore[AutoUser]) error {
		u := &AutoUser{Name: "Alice"}
		if err := tx.Insert(context.Background(), u); err != nil {
			return err
		}
		inserted = *u
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), inserted.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Tx_Rollback(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[AutoUser](db)
	boom := errors.New("something went wrong")

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO autousers \(name\) VALUES \(\?\) RETURNING id`).
		WithArgs("Bob").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(2)))
	mock.ExpectRollback()

	err := repo.Tx(context.Background(), func(tx oca.GenericStore[AutoUser]) error {
		u := &AutoUser{Name: "Bob"}
		if err := tx.Insert(context.Background(), u); err != nil {
			return err
		}
		return boom
	})

	assert.ErrorIs(t, err, boom)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Tx_BeginError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[AutoUser](db)

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	err := repo.Tx(context.Background(), func(tx oca.GenericStore[AutoUser]) error {
		return nil
	})

	assert.ErrorIs(t, err, sql.ErrConnDone)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Tx_MultipleOps(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[AutoUser](db)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO autousers \(name\) VALUES \(\?\) RETURNING id`).
		WithArgs("Carol").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(10)))
	mock.ExpectExec(`UPDATE autousers SET name = \? WHERE id = \?`).
		WithArgs("Carol Smith", int64(10)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Tx(context.Background(), func(tx oca.GenericStore[AutoUser]) error {
		u := &AutoUser{Name: "Carol"}
		if err := tx.Insert(context.Background(), u); err != nil {
			return err
		}
		u.Name = "Carol Smith"
		return tx.Update(context.Background(), u)
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Tx_WithDialect(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := oca.NewRepository[AutoUser](db, oca.WithDialect[AutoUser](query.PostgresDialect{}))

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO autousers \(name\) VALUES \(\$1\) RETURNING id`).
		WithArgs("Dave").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(7)))
	mock.ExpectCommit()

	err := repo.Tx(context.Background(), func(tx oca.GenericStore[AutoUser]) error {
		u := &AutoUser{Name: "Dave"}
		return tx.Insert(context.Background(), u)
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
