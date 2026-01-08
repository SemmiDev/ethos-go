package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
	"github.com/semmidev/ethos-go/internal/common/database"
)

type UserPostgresRepository struct {
	db database.DBTX
}

func NewUserPostgresRepository(db database.DBTX) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

func (r *UserPostgresRepository) Create(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (
			user_id, email, name, hashed_password, is_active,
			is_verified, verify_token, verify_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		u.UserID,
		u.Email,
		u.Name,
		u.HashedPassword,
		u.IsActive,
		u.IsVerified,
		u.VerifyToken,
		u.VerifyExpiresAt,
		u.PasswordResetToken,
		u.PasswordResetExpiresAt,
		u.CreatedAt,
		u.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return user.ErrAlreadyExists
			}
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserPostgresRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT
			user_id, email, name, hashed_password, is_active,
			is_verified, verify_token, verify_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var u user.User
	err := r.db.QueryRowxContext(ctx, query, email).StructScan(&u)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return &u, nil
}

func (r *UserPostgresRepository) FindByID(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	query := `
		SELECT
			user_id, email, name, hashed_password, is_active,
			is_verified, verify_token, verify_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		FROM users
		WHERE user_id = $1
	`

	var u user.User
	err := r.db.QueryRowxContext(ctx, query, userID).StructScan(&u)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return &u, nil
}

func (r *UserPostgresRepository) Update(ctx context.Context, u *user.User) error {
	u.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET
			email = $1,
			name = $2,
			hashed_password = $3,
			is_active = $4,
			is_verified = $5,
			verify_token = $6,
			verify_expires_at = $7,
			password_reset_token = $8,
			password_reset_expires_at = $9,
			updated_at = $10
		WHERE user_id = $11
	`

	res, err := r.db.ExecContext(ctx, query,
		u.Email,
		u.Name,
		u.HashedPassword,
		u.IsActive,
		u.IsVerified,
		u.VerifyToken,
		u.VerifyExpiresAt,
		u.PasswordResetToken,
		u.PasswordResetExpiresAt,
		u.UpdatedAt, // $10
		u.UserID,    // $11
	)

	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return user.ErrNotFound
	}
	return nil
}

func (r *UserPostgresRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM users WHERE user_id = $1`
	res, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return user.ErrNotFound
	}
	return nil
}
