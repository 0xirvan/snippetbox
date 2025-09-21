package mysql

import (
	"context"
	"database/sql"
	"errors"

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

var ErrNoUserFound = errors.New("no user found")

func (r *UserRepo) Save(ctx context.Context, u *user.User) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO users (email, password) VALUES (?, ?)",
		u.Email, u.Password,
	)
	return err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT * FROM users WHERE email = ?",
		email,
	)
	u := new(user.User)
	err := row.Scan(u.ID, u.Email, u.Password, u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoUserFound
	}
	return u, nil
}
