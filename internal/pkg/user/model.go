package user

import "errors"

type User struct {
	ID       string
	Username string
	Password string
}

var (
	ErrUserNotFound  = errors.New("user was not found")
	ErrUsedUsername  = errors.New("username already in use")
	ErrEmptyPassword = errors.New("empty password or username")
)
