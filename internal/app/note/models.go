package note

import "time"

type NoteResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Subject     string    `json:"subject"`
	Text        string    `json:"text"`
	TTL         *int64    `json:"ttl"`
	IsPublic    bool      `json:"isPublic"`
	PublicUsers *[]string `json:"publicUsers"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type PostRequest struct {
	UserID      string    `json:"userId"`
	Subject     string    `json:"subject"`
	Text        string    `json:"text"`
	TTL         *int64    `json:"ttl"`
	IsPublic    bool      `json:"isPublic"`
	PublicUsers *[]string `json:"publicUsers"`
}

type UpdateRequest struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Subject     string    `json:"subject"`
	Text        string    `json:"text"`
	TTL         *int64    `json:"ttl"`
	IsPublic    bool      `json:"isPublic"`
	PublicUsers *[]string `json:"publicUsers"`
}
