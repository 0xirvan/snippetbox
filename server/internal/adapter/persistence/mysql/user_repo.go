package mysql

import (
	"context"
	"database/sql"

	"github.com/0xirvan/snippetbox/internal/domain/user"
	"github.com/0xirvan/snippetbox/internal/port/auth"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) auth.AuthUserRepository {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Save(ctx context.Context, u *user.User) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO users (email, password) VALUES (?, ?)",
		u.Email, u.Password,
	)
	return err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	row, err := r.db.QueryContext(ctx,
		"SELECT id, email, password FROM users WHERE email = ?",
		email,
	)
	if err != nil {
		return nil, err
	}

	u := new(user.User)
	if err := row.Scan(u.ID, u.Email, u.Password); err != nil {
		return nil, err
	}
	return u, nil
}
