package database

import (
	"context"
)

// TransactionProvider provides transaction management for command handlers.
// Use this pattern when you need to share a transaction across multiple repositories.
//
// Usage:
//
//	type Adapters struct {
//	    UserRepo  UserRepository
//	    AuditRepo AuditLogRepository
//	}
//
//	func (p *TransactionProvider[Adapters]) Transact(ctx context.Context, fn func(adapters Adapters) error) error {
//	    return database.RunInTx(ctx, p.db, func(tx *sqlx.Tx) error {
//	        adapters := Adapters{
//	            UserRepo:  NewUserRepository(tx),
//	            AuditRepo: NewAuditLogRepository(tx),
//	        }
//	        return fn(adapters)
//	    })
//	}
//
// WARNING: This pattern should be used sparingly. In most cases, prefer the UpdateFn pattern
// where transactions are handled entirely within a single repository method.
type TransactionProvider struct {
	db DBTX
}

// NewTransactionProvider creates a new TransactionProvider.
func NewTransactionProvider(db DBTX) *TransactionProvider {
	return &TransactionProvider{db: db}
}

// DB returns the underlying database connection.
func (p *TransactionProvider) DB() DBTX {
	return p.db
}

// Transact executes a function within a transaction, providing a transaction handle.
// Use this when you need to work with multiple repositories in a single transaction.
//
// Example:
//
//	err := provider.Transact(ctx, func(tx database.DBTX) error {
//	    userRepo := adapters.NewUserRepository(tx)
//	    habitRepo := adapters.NewHabitRepository(tx)
//
//	    // Both operations run in the same transaction
//	    if err := userRepo.Create(ctx, user); err != nil {
//	        return err
//	    }
//	    return habitRepo.Create(ctx, habit)
//	})
func (p *TransactionProvider) Transact(ctx context.Context, fn func(tx DBTX) error) error {
	return RunInTx(ctx, p.db, fn)
}

// TransactWithResult executes a function within a transaction and returns a result.
func (p *TransactionProvider) TransactWithResult(ctx context.Context, fn func(tx DBTX) (any, error)) (any, error) {
	return RunInTxWithResult(ctx, p.db, fn)
}

// TransactWithOptions executes a function within a transaction with custom options.
func (p *TransactionProvider) TransactWithOptions(ctx context.Context, opts *TxOptions, fn func(tx DBTX) error) error {
	return RunInTxWithOptions(ctx, p.db, opts, fn)
}

// GenericTransactionProvider is a generic version that allows specifying the adapters type.
// This provides type-safe access to repository adapters within a transaction.
//
// Usage:
//
//	type Adapters struct {
//	    UserRepo  *UserRepository
//	    AuditRepo *AuditLogRepository
//	}
//
//	provider := NewGenericTransactionProvider(db, func(tx database.DBTX) Adapters {
//	    return Adapters{
//	        UserRepo:  NewUserRepository(tx),
//	        AuditRepo: NewAuditLogRepository(tx),
//	    }
//	})
//
//	err := provider.Transact(ctx, func(adapters Adapters) error {
//	    return adapters.UserRepo.Create(ctx, user)
//	})
type GenericTransactionProvider[T any] struct {
	db         DBTX
	adaptersFn func(tx DBTX) T
}

// NewGenericTransactionProvider creates a new GenericTransactionProvider.
// adaptersFn is a function that creates the adapters struct from a transaction.
func NewGenericTransactionProvider[T any](db DBTX, adaptersFn func(tx DBTX) T) *GenericTransactionProvider[T] {
	return &GenericTransactionProvider[T]{
		db:         db,
		adaptersFn: adaptersFn,
	}
}

// Transact executes a function within a transaction, providing typed adapters.
func (p *GenericTransactionProvider[T]) Transact(ctx context.Context, fn func(adapters T) error) error {
	return RunInTx(ctx, p.db, func(tx DBTX) error {
		adapters := p.adaptersFn(tx)
		return fn(adapters)
	})
}
