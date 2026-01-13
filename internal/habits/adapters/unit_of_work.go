package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/semmidev/ethos-go/internal/common/database"
	"github.com/semmidev/ethos-go/internal/habits/domain/habit"
	"github.com/semmidev/ethos-go/internal/habits/domain/habitlog"
)

// HabitsUnitOfWork coordinates transactions across habit-related repositories.
// It ensures that all repository operations within a transaction
// either succeed together or fail together.
//
// Usage without transaction (direct repository access):
//
//	habit, err := uow.Habits().GetHabit(ctx, habitID, userID)
//
// Usage with transaction:
//
//	err := uow.WithTransaction(ctx, func(txUow HabitsUnitOfWork) error {
//	    if err := txUow.HabitLogs().AddHabitLog(ctx, log); err != nil {
//	        return err
//	    }
//	    return txUow.Habits().UpsertStats(ctx, stats)
//	})
type HabitsUnitOfWork interface {
	// Habits returns the habit repository within this unit of work.
	Habits() habit.Repository

	// HabitLogs returns the habit log repository within this unit of work.
	HabitLogs() habitlog.Repository

	// WithTransaction executes a function within a transaction.
	// It automatically commits on success or rolls back on error/panic.
	// The callback receives a transactional UnitOfWork with repositories
	// that share the same transaction.
	WithTransaction(ctx context.Context, fn func(HabitsUnitOfWork) error) error
}

// habitsUnitOfWork is the PostgreSQL implementation of HabitsUnitOfWork.
type habitsUnitOfWork struct {
	db            database.DBTX
	habitRepo     habit.Repository
	logRepo       habitlog.Repository
	inTransaction bool
}

// NewHabitsUnitOfWork creates a new habits unit of work.
func NewHabitsUnitOfWork(db database.DBTX) HabitsUnitOfWork {
	return &habitsUnitOfWork{
		db:        db,
		habitRepo: NewHabitPostgresRepository(db),
		logRepo:   NewHabitLogPostgresRepository(db),
	}
}

// Habits returns the habit repository.
// When in a transaction, it returns the transactional repository.
func (uow *habitsUnitOfWork) Habits() habit.Repository {
	return uow.habitRepo
}

// HabitLogs returns the habit log repository.
// When in a transaction, it returns the transactional repository.
func (uow *habitsUnitOfWork) HabitLogs() habitlog.Repository {
	return uow.logRepo
}

// WithTransaction executes a function within a transaction.
// This is the recommended way to use transactions as it handles
// commit and rollback automatically, including panic recovery.
func (uow *habitsUnitOfWork) WithTransaction(ctx context.Context, fn func(HabitsUnitOfWork) error) (err error) {
	// If already in a transaction, just run the function (nested transaction support)
	if uow.inTransaction {
		return fn(uow)
	}

	// Get the underlying database connection to begin a transaction
	db := uow.db

	// Unwrap TracedDBTX if necessary to access BeginTxx
	if traced, ok := db.(*database.TracedDBTX); ok {
		db = traced.Unwrap()
	}

	// If it's already a transaction, just run the function
	if tx, ok := db.(*sqlx.Tx); ok {
		txUow := &habitsUnitOfWork{
			db:            tx,
			habitRepo:     NewHabitPostgresRepository(tx),
			logRepo:       NewHabitLogPostgresRepository(tx),
			inTransaction: true,
		}
		return fn(txUow)
	}

	// It must be *sqlx.DB to start a transaction
	conn, ok := db.(*sqlx.DB)
	if !ok {
		return errors.New("WithTransaction: db must be *sqlx.DB or *sqlx.Tx")
	}

	tx, err := conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Handle panic recovery
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// Create a new UnitOfWork with transactional repositories
	txUow := &habitsUnitOfWork{
		db:            tx,
		habitRepo:     NewHabitPostgresRepository(tx),
		logRepo:       NewHabitLogPostgresRepository(tx),
		inTransaction: true,
	}

	err = fn(txUow)
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
