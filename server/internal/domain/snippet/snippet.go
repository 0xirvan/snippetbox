package snippet

import "time"

type Snippet struct {
	ID      uint
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
