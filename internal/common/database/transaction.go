package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// RunInTx executes a function within a database transaction.
// It automatically commits the transaction if the function returns nil,
// or rolls back if the function returns an error or panics.
//
// Usage:
//
//	err := database.RunInTx(ctx, db, func(tx *sqlx.Tx) error {
//	    // perform database operations
//	    return nil
//	})
func RunInTx(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) error) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic
			_ = tx.Rollback()
			panic(p) // Re-throw panic after rollback
		}
	}()

	err = fn(tx)
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
//	user, err := database.RunInTxWithResult(ctx, db, func(tx *sqlx.Tx) (*User, error) {
//	    // perform database operations
//	    return user, nil
//	})
func RunInTxWithResult[T any](ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) (T, error)) (result T, err error) {
	tx, err := db.BeginTxx(ctx, nil)
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

	result, err = fn(tx)
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
func RunInTxWithOptions(ctx context.Context, db *sqlx.DB, opts *TxOptions, fn func(tx *sqlx.Tx) error) (err error) {
	tx, err := db.BeginTxx(ctx, opts.toSQLTxOptions())
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	err = fn(tx)
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
