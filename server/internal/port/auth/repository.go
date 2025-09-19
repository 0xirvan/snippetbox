package auth

import (
	"context"

	"github.com/0xirvan/snippetbox/internal/domain/user"
)

type AuthUserRepository interface {
	Save(ctx context.Context, u *user.User) error
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}
