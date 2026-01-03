package database

import (
	"context"
	"database/sql"
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

	err = RunInTx(context.Background(), sqlxDB, func(tx *sqlx.Tx) error {
		_, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "test")
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

	err = RunInTx(context.Background(), sqlxDB, func(tx *sqlx.Tx) error {
		_, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "test")
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

	_ = RunInTx(context.Background(), sqlxDB, func(tx *sqlx.Tx) error {
		panic("something went wrong")
	})
}

func TestRunInTxWithResult_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id FROM users").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))
	mock.ExpectCommit()

	result, err := RunInTxWithResult(context.Background(), sqlxDB, func(tx *sqlx.Tx) (int, error) {
		var id int
		err := tx.QueryRow("SELECT id FROM users WHERE name = ?", "test").Scan(&id)
		return id, err
	})

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if result != 42 {
		t.Errorf("expected result 42, got: %d", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRunInTxWithResult_RollbackOnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id FROM users").WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	result, err := RunInTxWithResult(context.Background(), sqlxDB, func(tx *sqlx.Tx) (int, error) {
		var id int
		err := tx.QueryRow("SELECT id FROM users WHERE name = ?", "test").Scan(&id)
		return id, err
	})

	if err == nil {
		t.Error("expected error, got nil")
	}

	if result != 0 {
		t.Errorf("expected zero value result, got: %d", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestTransactionProvider_Transact(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	provider := NewTransactionProvider(sqlxDB)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = provider.Transact(context.Background(), func(tx *sqlx.Tx) error {
		_, err := tx.Exec("UPDATE users SET name = ? WHERE id = ?", "updated", 1)
		return err
	})

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGenericTransactionProvider_Transact(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	// Define adapters struct
	type TestAdapters struct {
		tx *sqlx.Tx
	}

	provider := NewGenericTransactionProvider(sqlxDB, func(tx *sqlx.Tx) TestAdapters {
		return TestAdapters{tx: tx}
	})

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO audit_log").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = provider.Transact(context.Background(), func(adapters TestAdapters) error {
		_, err := adapters.tx.Exec("INSERT INTO audit_log (message) VALUES (?)", "test")
		return err
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
