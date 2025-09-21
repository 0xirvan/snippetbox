package session

import "time"

type SessionEntry struct {
	ID     uint
	UserID uint
	Expiry time.Time
}
