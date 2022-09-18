package note

import (
	"errors"
	"time"
)

type Note struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Subject     string    `json:"subject"`
	Text        string    `json:"text"`
	TTL         *int      `json:"ttl"`
	IsPublic    bool      `json:"isPublic"`
	PublicUsers *[]string `json:"publicUsers"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var (
	ErrEmptyNote    = errors.New("empty note text")
	ErrNoteNotFound = errors.New("note not found")
)
