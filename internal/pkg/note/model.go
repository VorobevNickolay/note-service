package note

import (
	"errors"
	"time"
)

type Note struct {
	ID          string   `json:"id"`
	UserID      string   `json:"userId"`
	Text        string   `json:"text"`
	TTL         string   `json:"ttl"`
	IsPublic    bool     `json:"isPublic"`
	PublicUsers []string `json:"publicUsers"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// todo: add parameters
}

var (
	ErrEmptyNote    = errors.New("empty note text")
	ErrNoteNotFound = errors.New("note not found")
)
