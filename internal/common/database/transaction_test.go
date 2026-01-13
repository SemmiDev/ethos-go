package database

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestRunInTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = RunInTx(context.Background(), sqlxDB, func(tx DBTX) error {
		_, err := tx.ExecContext(context.Background(), "INSERT INTO users (name) VALUES (?)", "test")
		return err
	})

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRunInTx_RollbackOnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	expectedErr := errors.New("insert failed")

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users").WillReturnError(expectedErr)
	mock.ExpectRollback()

	err = RunInTx(context.Background(), sqlxDB, func(tx DBTX) error {
		_, err := tx.ExecContext(context.Background(), "INSERT INTO users (name) VALUES (?)", "test")
		return err
	})

	if err == nil {
		t.Error("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error to wrap %v, got: %v", expectedErr, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRunInTx_RollbackOnPanic(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectRollback()

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, got none")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	}()

	_ = RunInTx(context.Background(), sqlxDB, func(tx DBTX) error {
		panic("something went wrong")
	})
}

func TestRunInTx_NestedTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO logs").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test nested transaction - inner RunInTx should just run the function
	err = RunInTx(context.Background(), sqlxDB, func(tx DBTX) error {
		_, err := tx.ExecContext(context.Background(), "INSERT INTO users (name) VALUES (?)", "test")
		if err != nil {
			return err
		}

		// Nested transaction - should not start a new transaction
		return RunInTx(context.Background(), tx, func(innerTx DBTX) error {
			_, err := innerTx.ExecContext(context.Background(), "INSERT INTO logs (message) VALUES (?)", "log")
			return err
		})
	})

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDBTX_Interface(t *testing.T) {
	// Compile-time checks are in db_interface.go, but this test
	// documents that both types implement the interface
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	// Ensure *sqlx.DB implements DBTX
	var _ DBTX = sqlxDB

	t.Log("Both *sqlx.DB and *sqlx.Tx implement DBTX interface")
}
