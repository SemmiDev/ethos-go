package service

import "context"

// PasswordHasher provides secure password hashing and verification.
// We define this as an interface so the app layer doesn't depend on
// a specific hashing algorithm (bcrypt, argon2, etc). The infrastructure
// layer can switch implementations without changing app code.
type PasswordHasher interface {
	// Hash converts a plain text password into a secure hash.
	// The hash can be safely stored in the database.
	Hash(ctx context.Context, password string) (string, error)

	// Compare checks if a plain text password matches a hash.
	// Returns true if they match, false otherwise.
	Compare(ctx context.Context, hashedPassword, plainPassword string) (bool, error)
}
