package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// RunInTx executes a function within a database transaction.
// It automatically commits on success or rolls back on error/panic.
//
// The db argument must be *sqlx.DB or *TracedDBTX wrapping *sqlx.DB to start a transaction.
// If *sqlx.Tx or other DBTX is passed, it executes without starting a new transaction
// (nested transaction support).
//
// Usage:
//
//	err := database.RunInTx(ctx, db, func(tx database.DBTX) error {
//	    // All operations here run within the same transaction
//	    if err := repo1.Create(ctx, entity); err != nil {
//	        return err // triggers rollback
//	    }
//	    return repo2.Update(ctx, entity) // success commits
//	})
func RunInTx(ctx context.Context, db DBTX, fn func(tx DBTX) error) (err error) {
	var wrapTx func(DBTX) DBTX

	// Unwrap TracedDBTX to get access to the underlying *sqlx.DB for BeginTxx,
	// but remember to wrap the transaction later for tracing.
	if traced, ok := db.(*TracedDBTX); ok {
		db = traced.Unwrap()
		wrapTx = func(tx DBTX) DBTX {
			return NewTracedDBTX(tx)
		}
	}

	// If it's already a transaction, just run the function (nested transaction support)
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
