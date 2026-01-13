package adapters

import (
	"context"

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
	tx            *sqlx.Tx // nil when not in a transaction
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
func (uow *habitsUnitOfWork) WithTransaction(ctx context.Context, fn func(HabitsUnitOfWork) error) error {
	// If already in a transaction, just run the function (nested transaction support)
	if uow.inTransaction {
		return fn(uow)
	}

	// Use the database package's RunInTx for proper transaction handling
	return database.RunInTx(ctx, uow.db, func(tx database.DBTX) error {
		// Create a new UnitOfWork with transactional repositories
		txUow := &habitsUnitOfWork{
			db:            tx,
			habitRepo:     NewHabitPostgresRepository(tx),
			logRepo:       NewHabitLogPostgresRepository(tx),
			inTransaction: true,
		}
		return fn(txUow)
	})
}
