package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// DBTX is an interface that both *sqlx.DB and *sqlx.Tx implement.
// This allows repositories to work with either a direct connection or a transaction.
//
// Usage:
//
//	type UserRepository struct {
//	    db database.DBTX
//	}
//
//	func NewUserRepository(db database.DBTX) *UserRepository {
//	    return &UserRepository{db: db}
//	}
//
// This pattern allows the repository to be used directly with a database connection
// or within a transaction by passing a *sqlx.Tx instead of *sqlx.DB.
type DBTX interface {
	// ExecContext executes a query without returning any rows.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)

	// GetContext queries the database and scans a single row into dest.
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error

	// SelectContext queries the database and scans multiple rows into dest.
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error

	// QueryRowxContext queries the database and returns a single Row.
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row

	// QueryxContext queries the database and returns an *sqlx.Rows.
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)

	// PreparexContext prepares a statement.
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)

	// Rebind transforms a query from QUESTION to the DB driver's bind type.
	Rebind(query string) string

	// DriverName returns the driverName used by this DB.
	DriverName() string
}

// TxOptions represents transaction options for custom transaction configurations.
type TxOptions struct {
	// Isolation is the transaction isolation level.
	// If zero, the driver or database's default level is used.
	Isolation IsolationLevel

	// ReadOnly indicates whether the transaction is read-only.
	ReadOnly bool
}

// IsolationLevel is the transaction isolation level used in TxOptions.
type IsolationLevel int

// Various isolation levels that drivers may support.
const (
	LevelDefault IsolationLevel = iota
	LevelReadUncommitted
	LevelReadCommitted
	LevelWriteCommitted
	LevelRepeatableRead
	LevelSnapshot
	LevelSerializable
	LevelLinearizable
)

// toSQLTxOptions converts TxOptions to *sql.TxOptions.
func (o *TxOptions) toSQLTxOptions() *sql.TxOptions {
	if o == nil {
		return nil
	}
	return &sql.TxOptions{
		Isolation: sql.IsolationLevel(o.Isolation),
		ReadOnly:  o.ReadOnly,
	}
}

// Compile-time checks to ensure *sqlx.DB and *sqlx.Tx implement DBTX.
var (
	_ DBTX = (*sqlx.DB)(nil)
	_ DBTX = (*sqlx.Tx)(nil)
)
