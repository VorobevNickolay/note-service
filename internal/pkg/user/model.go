package user

import "errors"

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	ErrUserNotFound  = errors.New("user was not found")
	ErrUsedUsername  = errors.New("username already in use")
	ErrEmptyPassword = errors.New("empty password or username")
)
