package note

type NoteResponse struct {
}

type PostRequest struct {
	UserID      string    `json:"userId"`
	Subject     string    `json:"subject"`
	Text        string    `json:"text"`
	TTL         *int      `json:"ttl"`
	IsPublic    bool      `json:"isPublic"`
	PublicUsers *[]string `json:"publicUsers"`
}

type UpdateRequest struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Subject     string    `json:"subject"`
	Text        string    `json:"text"`
	TTL         *int      `json:"ttl"`
	IsPublic    bool      `json:"isPublic"`
	PublicUsers *[]string `json:"publicUsers"`
}
