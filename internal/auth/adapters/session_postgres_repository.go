package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/semmidev/ethos-go/internal/auth/domain/session"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
	"github.com/semmidev/ethos-go/internal/common/model"
)

// ListSessions retrieves a paginated list of sessions for a user with filtering options.
func (r *SessionPostgresRepository) ListSessions(
	ctx context.Context,
	userID uuid.UUID,
	includeBlocked, includeExpired bool,
	filter model.Filter,
) ([]*session.Session, int, error) {
	// Base query
	query := `
		SELECT
			session_id, user_id, refresh_token, user_agent,
			client_ip, is_blocked, expires_at, created_at, updated_at
		FROM sessions
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argIdx := 2

	// Conditions
	if !includeBlocked {
		query += fmt.Sprintf(" AND is_blocked = false")
	}
	if !includeExpired {
		query += fmt.Sprintf(" AND expires_at > NOW()")
	}

	// Dynamic sorting
	orderColumn := "created_at"
	orderDirection := "DESC"

	if filter.SortBy != "" {
		// Whitelist allowed columns to prevent SQL injection
		allowedColumns := map[string]bool{
			"created_at": true,
			"expires_at": true,
			"is_blocked": true,
		}
		if allowedColumns[filter.SortBy] {
			orderColumn = filter.SortBy
		}
	}
	if filter.SortDirection == "asc" {
		orderDirection = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", orderColumn, orderDirection)

	// Pagination
	limit := filter.PerPage
	offset := filter.GetOffset()

	// Get total count (for pagination)
	countQuery := "SELECT COUNT(*) FROM (" + query + ") as count_table"
	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, r.translateError(err, "count sessions")
	}

	// Apply limit/offset to main query
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	var sessions []*session.Session
	err = r.db.SelectContext(ctx, &sessions, query, args...)
	if err != nil {
		return nil, 0, r.translateError(err, "list sessions")
	}

	return sessions, totalCount, nil
}

// SessionPostgresRepository implements the SessionRepository interface.
type SessionPostgresRepository struct {
	db *sqlx.DB
}

func NewSessionPostgresRepository(db *sqlx.DB) *SessionPostgresRepository {
	return &SessionPostgresRepository{
		db: db,
	}
}

// Create inserts a new session into the database.
func (r *SessionPostgresRepository) Create(ctx context.Context, s *session.Session) error {
	query := `
		INSERT INTO sessions (
			session_id, user_id, refresh_token, user_agent,
			client_ip, is_blocked, expires_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		s.SessionID,
		s.UserID,
		s.RefreshToken,
		s.UserAgent,
		s.ClientIP,
		s.IsBlocked,
		s.ExpiresAt,
		s.CreatedAt,
		s.UpdatedAt,
	)

	if err != nil {
		return r.translateError(err, "create session")
	}

	return nil
}

// FindByID retrieves a session by its unique identifier.
func (r *SessionPostgresRepository) FindByID(ctx context.Context, sessionID uuid.UUID) (*session.Session, error) {
	query := `
		SELECT
			session_id, user_id, refresh_token, user_agent,
			client_ip, is_blocked, expires_at, created_at, updated_at
		FROM sessions
		WHERE session_id = $1
	`

	var s session.Session
	err := r.db.QueryRowxContext(ctx, query, sessionID).StructScan(&s)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, session.ErrNotFound
		}
		return nil, r.translateError(err, "find session by id")
	}

	return &s, nil
}

// FindByRefreshToken looks up a session using its refresh token.
func (r *SessionPostgresRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (*session.Session, error) {
	query := `
		SELECT
			session_id, user_id, refresh_token, user_agent,
			client_ip, is_blocked, expires_at, created_at, updated_at
		FROM sessions
		WHERE refresh_token = $1
	`

	var s session.Session
	err := r.db.QueryRowxContext(ctx, query, refreshToken).StructScan(&s)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, session.ErrNotFound
		}
		return nil, r.translateError(err, "find session by refresh token")
	}

	return &s, nil
}

// FindAllByUserID returns all sessions for a specific user.
func (r *SessionPostgresRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*session.Session, error) {
	query := `
		SELECT
			session_id, user_id, refresh_token, user_agent,
			client_ip, is_blocked, expires_at, created_at, updated_at
		FROM sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var sessions []*session.Session
	err := r.db.SelectContext(ctx, &sessions, query, userID)
	if err != nil {
		return nil, r.translateError(err, "find sessions by user id")
	}

	return sessions, nil
}

// Update modifies an existing session in the database.
func (r *SessionPostgresRepository) Update(ctx context.Context, s *session.Session) error {
	query := `
		UPDATE sessions
		SET
			refresh_token = $2,
			is_blocked = $3,
			expires_at = $4,
			updated_at = $5
		WHERE session_id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		s.SessionID,
		s.RefreshToken,
		s.IsBlocked,
		s.ExpiresAt,
		s.UpdatedAt,
	)

	if err != nil {
		return r.translateError(err, "update session")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return session.ErrNotFound
	}

	return nil
}

// Delete permanently removes a session from the database.
func (r *SessionPostgresRepository) Delete(ctx context.Context, sessionID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE session_id = $1`

	result, err := r.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		return r.translateError(err, "delete session")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return session.ErrNotFound
	}

	return nil
}

// DeleteAllByUserID removes all sessions for a user.
func (r *SessionPostgresRepository) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return r.translateError(err, "delete user sessions")
	}

	return nil
}

// DeleteExpired removes all sessions that have passed their expiration time.
func (r *SessionPostgresRepository) DeleteExpired(ctx context.Context) (int64, error) {
	query := `
		DELETE FROM sessions
		WHERE expires_at < NOW() AND NOT is_blocked
	`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, r.translateError(err, "delete expired sessions")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("check rows affected: %w", err)
	}

	return rowsAffected, nil
}

// translateError converts database-specific errors to domain errors.
func (r *SessionPostgresRepository) translateError(err error, operation string) error {
	// Check for PostgreSQL-specific errors
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return fmt.Errorf("session already exists: %w", err)

		case "23503": // foreign_key_violation
			return fmt.Errorf("%w: user not found", user.ErrNotFound)

		case "23502": // not_null_violation
			return fmt.Errorf("required field missing: %s: %w", pgErr.Column, err)

		case "22P02": // invalid_text_representation (bad UUID format)
			return fmt.Errorf("%w: invalid id format", errors.New("invalid UUID"))

		case "42P01": // undefined_table
			return fmt.Errorf("database schema error: table not found: %w", err)
		}
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("database operation timed out during %s: %w", operation, err)
	}

	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("database operation canceled during %s: %w", operation, err)
	}

	return fmt.Errorf("database error during %s: %w", operation, err)
}
