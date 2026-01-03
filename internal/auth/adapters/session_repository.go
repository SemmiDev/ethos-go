package adapters

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/auth/domain/session"
)

// InMemorySessionRepository is a simple in-memory implementation for development/testing
type InMemorySessionRepository struct {
	mu       sync.RWMutex
	sessions map[uuid.UUID]*session.Session
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{
		sessions: make(map[uuid.UUID]*session.Session),
	}
}

func (r *InMemorySessionRepository) Create(ctx context.Context, sess *session.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[sess.SessionID] = sess
	return nil
}

func (r *InMemorySessionRepository) FindByID(ctx context.Context, sessionID uuid.UUID) (*session.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sess, exists := r.sessions[sessionID]
	if !exists {
		return nil, session.ErrNotFound
	}

	return sess, nil
}

func (r *InMemorySessionRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (*session.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, session := range r.sessions {
		if session.RefreshToken == refreshToken {
			return session, nil
		}
	}

	return nil, session.ErrNotFound
}

func (r *InMemorySessionRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*session.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var sessions []*session.Session
	for _, session := range r.sessions {
		if session.UserID == userID {
			sessions = append(sessions, session)
		}
	}

	return sessions, nil
}

func (r *InMemorySessionRepository) Update(ctx context.Context, sess *session.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[sess.SessionID]; !exists {
		return session.ErrNotFound
	}

	r.sessions[sess.SessionID] = sess
	return nil
}

func (r *InMemorySessionRepository) Delete(ctx context.Context, sessionID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[sessionID]; !exists {
		return session.ErrNotFound
	}

	delete(r.sessions, sessionID)
	return nil
}

func (r *InMemorySessionRepository) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for sessionID, session := range r.sessions {
		if session.UserID == userID {
			delete(r.sessions, sessionID)
		}
	}

	return nil
}

func (r *InMemorySessionRepository) DeleteExpired(ctx context.Context) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var count int64
	now := time.Now()

	for sessionID, session := range r.sessions {
		if session.ExpiresAt.Before(now) {
			delete(r.sessions, sessionID)
			count++
		}
	}

	return count, nil
}
