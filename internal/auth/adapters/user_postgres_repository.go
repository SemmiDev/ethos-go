package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
	// Convert domain entity to database model
	model := UserModelFromUser(u)

	query := `
		INSERT INTO users (
			user_id, email, name, hashed_password, auth_provider, auth_provider_id,
			timezone, is_active, is_verified, verify_token, verify_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.db.ExecContext(ctx, query,
		model.UserID,
		model.Email,
		model.Name,
		model.HashedPassword,
		model.AuthProvider,
		model.AuthProviderID,
		model.Timezone,
		model.IsActive,
		model.IsVerified,
		model.VerifyToken,
		model.VerifyExpiresAt,
		model.PasswordResetToken,
		model.PasswordResetExpiresAt,
		model.CreatedAt,
		model.UpdatedAt,
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
			user_id, email, name, hashed_password, auth_provider, auth_provider_id,
			timezone, is_active, is_verified, verify_token, verify_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var model UserModel
	err := r.db.QueryRowxContext(ctx, query, email).StructScan(&model)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return model.ToUser(), nil
}

func (r *UserPostgresRepository) FindByID(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	query := `
		SELECT
			user_id, email, name, hashed_password, auth_provider, auth_provider_id,
			timezone, is_active, is_verified, verify_token, verify_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		FROM users
		WHERE user_id = $1
	`

	var model UserModel
	err := r.db.QueryRowxContext(ctx, query, userID).StructScan(&model)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return model.ToUser(), nil
}

func (r *UserPostgresRepository) FindByAuthProvider(ctx context.Context, provider, providerID string) (*user.User, error) {
	query := `
		SELECT
			user_id, email, name, hashed_password, auth_provider, auth_provider_id,
			timezone, is_active, is_verified, verify_token, verify_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		FROM users
		WHERE auth_provider = $1 AND auth_provider_id = $2
	`

	var model UserModel
	err := r.db.QueryRowxContext(ctx, query, provider, providerID).StructScan(&model)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrNotFound
		}
		return nil, fmt.Errorf("find user by auth provider: %w", err)
	}

	return model.ToUser(), nil
}

func (r *UserPostgresRepository) Update(ctx context.Context, u *user.User) error {
	// Convert domain entity to database model
	model := UserModelFromUser(u)

	query := `
		UPDATE users
		SET
			email = $1,
			name = $2,
			hashed_password = $3,
			auth_provider = $4,
			auth_provider_id = $5,
			timezone = $6,
			is_active = $7,
			is_verified = $8,
			verify_token = $9,
			verify_expires_at = $10,
			password_reset_token = $11,
			password_reset_expires_at = $12,
			updated_at = $13
		WHERE user_id = $14
	`

	res, err := r.db.ExecContext(ctx, query,
		model.Email,
		model.Name,
		model.HashedPassword,
		model.AuthProvider,
		model.AuthProviderID,
		model.Timezone,
		model.IsActive,
		model.IsVerified,
		model.VerifyToken,
		model.VerifyExpiresAt,
		model.PasswordResetToken,
		model.PasswordResetExpiresAt,
		model.UpdatedAt,
		model.UserID,
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
