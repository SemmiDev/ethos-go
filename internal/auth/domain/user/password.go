package user

// HashedPassword is a value object representing a hashed password.
// This provides type safety and prevents accidental use of raw passwords.
type HashedPassword struct {
	value string
}

// NewHashedPassword creates a new HashedPassword value object.
// Note: This does NOT perform hashing - it wraps an already-hashed password.
// Use the PasswordHasher service interface to hash raw passwords.
func NewHashedPassword(hashedValue string) HashedPassword {
	return HashedPassword{value: hashedValue}
}

// UnmarshalHashedPasswordFromDatabase reconstructs a HashedPassword from database storage.
func UnmarshalHashedPasswordFromDatabase(value string) HashedPassword {
	return HashedPassword{value: value}
}

// String returns the hashed password string.
// This is safe to store in the database but should never be shown to users.
func (p HashedPassword) String() string {
	return p.value
}

// IsEmpty returns true if the password is empty.
func (p HashedPassword) IsEmpty() bool {
	return p.value == ""
}

// Equals checks if two hashed passwords are the same.
func (p HashedPassword) Equals(other HashedPassword) bool {
	return p.value == other.value
}
