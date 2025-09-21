package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	sessionDomain "github.com/0xirvan/snippetbox/internal/domain/session"
	sessionPort "github.com/0xirvan/snippetbox/internal/port/session"
)

type InMemSessionRepo struct {
	mu       sync.RWMutex
	sessions map[string]sessionDomain.SessionEntry
}

var (
	ErrNoUserFound    = errors.New("no user found")
	ErrSessionExpired = errors.New("expired")
)

func NewInMemSessionRepo() sessionPort.SessionRepository {
	return &InMemSessionRepo{
		sessions: make(map[string]sessionDomain.SessionEntry),
	}
}

func (r *InMemSessionRepo) Save(ctx context.Context, sessionID string, userID uint, expiry time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[sessionID] = sessionDomain.SessionEntry{
		UserID: userID,
		Expiry: time.Now().Add(expiry),
	}
	return nil
}

func (r *InMemSessionRepo) Delete(ctx context.Context, sessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, sessionID)
	return nil
}

func (r *InMemSessionRepo) Get(ctx context.Context, sessionID string) (userID uint, err error) {
	r.mu.RLock()
	s, ok := r.sessions[sessionID]
	r.mu.RUnlock()
	if !ok {
		return 0, ErrNoUserFound
	}
	if time.Now().After(s.Expiry) {
		r.Delete(ctx, sessionID)
		return 0, ErrSessionExpired
	}
	return s.UserID, nil
}
