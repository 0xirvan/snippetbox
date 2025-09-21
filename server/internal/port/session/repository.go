package session

import (
	"context"
	"time"
)

type SessionRepository interface {
	Save(ctx context.Context, sessionID string, userID uint, expiry time.Duration) error
	Get(ctx context.Context, sessionID string) (userID uint, err error)
	Delete(ctx context.Context, sessionID string) error
}
