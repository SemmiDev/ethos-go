package adapters

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
)

var (
	ErrUserNotFound      = user.ErrNotFound
	ErrUserAlreadyExists = user.ErrAlreadyExists
)

// InMemoryUserRepository is a simple in-memory implementation for development/testing
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*user.User // keyed by email
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*user.User),
	}
}

func (r *InMemoryUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, exists := r.users[email]
	if !exists {
		return nil, ErrUserNotFound
	}

	return u, nil
}

func (r *InMemoryUserRepository) FindByID(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.UserID == userID {
			return u, nil
		}
	}

	return nil, ErrUserNotFound
}

func (r *InMemoryUserRepository) Create(ctx context.Context, u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[u.Email]; exists {
		return ErrUserAlreadyExists
	}

	r.users[u.Email] = u
	return nil
}

func (r *InMemoryUserRepository) Update(ctx context.Context, u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[u.Email]; !exists {
		return ErrUserNotFound
	}

	r.users[u.Email] = u
	return nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for email, u := range r.users {
		if u.UserID == userID {
			delete(r.users, email)
			return nil
		}
	}

	return ErrUserNotFound
}
