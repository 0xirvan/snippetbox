package user

import "time"

type User struct {
	ID        uint
	Email     string
	Password  string
	CreatedAt time.Time
}

func NewUser(email, password string) User {
	return User{
		Email:    email,
		Password: password,
	}
}
