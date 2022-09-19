package note

import (
	"errors"
	"time"
)

type Note struct {
	ID          string
	UserID      string
	Subject     string
	Text        string
	TTL         *int64
	IsPublic    bool
	PublicUsers *[]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var (
	ErrEmptyNote    = errors.New("empty note text")
	ErrNoteNotFound = errors.New("note not found")
)
