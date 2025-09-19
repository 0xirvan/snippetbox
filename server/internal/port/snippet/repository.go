package snippet

import (
	"context"

	"github.com/0xirvan/snippetbox/internal/domain/snippet"
)

type SnippetRepository interface {
	Save(ctx context.Context, s *snippet.Snippet) error
	GetByID(ctx context.Context, id uint) (*snippet.Snippet, error)
	Latest(ctx context.Context) ([]*snippet.Snippet, error) // Latest 10
}
