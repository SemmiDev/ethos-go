package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// RunInTx executes a function within a database transaction.
// The db argument must be *sqlx.DB or *TracedDBTX wrapping *sqlx.DB to start a transaction.
// If *sqlx.Tx or other DBTX is passed, it executes without starting a new transaction.
func RunInTx(ctx context.Context, db DBTX, fn func(tx DBTX) error) (err error) {
	var wrapTx func(DBTX) DBTX

	// Unwrap TracedDBTX to get redundant access to the underlying *sqlx.DB for BeginTxx,
	// but remember to wrap the transaction later for tracing.
	if traced, ok := db.(*TracedDBTX); ok {
		db = traced.Unwrap()
		wrapTx = func(tx DBTX) DBTX {
			return NewTracedDBTX(tx)
		}
	}

	// If it's already a transaction, just run the function
	if tx, ok := db.(*sqlx.Tx); ok {
		var txArg DBTX = tx
		if wrapTx != nil {
			txArg = wrapTx(tx)
		}
		return fn(txArg)
	}

	// It must be *sqlx.DB to start a transaction
	conn, ok := db.(*sqlx.DB)
	if !ok {
		return errors.New("RunInTx: db must be *sqlx.DB, *sqlx.Tx or *TracedDBTX")
	}

	tx, err := conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	var txArg DBTX = tx
	if wrapTx != nil {
		txArg = wrapTx(tx)
	}

	err = fn(txArg)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Join(err, fmt.Errorf("rollback: %w", rbErr))
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// RunInTxWithResult executes a function within a database transaction and returns a result.
// It automatically commits the transaction if the function returns nil error,
// or rolls back if the function returns an error or panics.
//
// Usage:
//
//	user, err := database.RunInTxWithResult(ctx, db, func(tx database.DBTX) (*User, error) {
//	    // perform database operations
//	    return user, nil
//	})
func RunInTxWithResult[T any](ctx context.Context, db DBTX, fn func(tx DBTX) (T, error)) (result T, err error) {
	var wrapTx func(DBTX) DBTX

	if traced, ok := db.(*TracedDBTX); ok {
		db = traced.Unwrap()
		wrapTx = func(tx DBTX) DBTX {
			return NewTracedDBTX(tx)
		}
	}

	// Just in case it's already a tx (though original only allowed *sqlx.DB)
	if tx, ok := db.(*sqlx.Tx); ok {
		var txArg DBTX = tx
		if wrapTx != nil {
			txArg = wrapTx(tx)
		}
		return fn(txArg)
	}

	conn, ok := db.(*sqlx.DB)
	if !ok {
		return result, errors.New("RunInTxWithResult: db must be *sqlx.DB, *sqlx.Tx or *TracedDBTX")
	}

	tx, err := conn.BeginTxx(ctx, nil)
	if err != nil {
		return result, fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic
			_ = tx.Rollback()
			panic(p) // Re-throw panic after rollback
		}
	}()

	var txArg DBTX = tx
	if wrapTx != nil {
		txArg = wrapTx(tx)
	}

	result, err = fn(txArg)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return result, errors.Join(err, fmt.Errorf("rollback: %w", rbErr))
		}
		return result, err
	}

	if err = tx.Commit(); err != nil {
		return result, fmt.Errorf("commit: %w", err)
	}

	return result, nil
}

// RunInTxWithOptions executes a function within a database transaction with custom options.
// This allows specifying isolation levels and read-only transactions.
//
// Usage:
//
//	opts := &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true}
//	err := database.RunInTxWithOptions(ctx, db, opts, func(tx *sqlx.Tx) error {
//	    // perform database operations
//	    return nil
//	})
//
// RunInTxWithOptions executes a function within a database transaction with custom options.
// This allows specifying isolation levels and read-only transactions.
//
// Usage:
//
//	opts := &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true}
//	err := database.RunInTxWithOptions(ctx, db, opts, func(tx database.DBTX) error {
//	    // perform database operations
//	    return nil
//	})
func RunInTxWithOptions(ctx context.Context, db DBTX, opts *TxOptions, fn func(tx DBTX) error) (err error) {
	var wrapTx func(DBTX) DBTX

	if traced, ok := db.(*TracedDBTX); ok {
		db = traced.Unwrap()
		wrapTx = func(tx DBTX) DBTX {
			return NewTracedDBTX(tx)
		}
	}

	// Just in case it's already a tx
	if tx, ok := db.(*sqlx.Tx); ok {
		var txArg DBTX = tx
		if wrapTx != nil {
			txArg = wrapTx(tx)
		}
		return fn(txArg)
	}

	conn, ok := db.(*sqlx.DB)
	if !ok {
		return errors.New("RunInTxWithOptions: db must be *sqlx.DB, *sqlx.Tx or *TracedDBTX")
	}

	tx, err := conn.BeginTxx(ctx, opts.toSQLTxOptions())
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	var txArg DBTX = tx
	if wrapTx != nil {
		txArg = wrapTx(tx)
	}

	err = fn(txArg)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Join(err, fmt.Errorf("rollback: %w", rbErr))
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
