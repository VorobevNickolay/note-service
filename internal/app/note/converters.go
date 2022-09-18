package note

import notepkg "note-service/internal/pkg/note"

func updateRequestToNote(request UpdateRequest) notepkg.Note {
	return notepkg.Note{
		ID:          request.ID,
		UserID:      request.UserID,
		Subject:     request.Subject,
		Text:        request.Text,
		TTL:         request.TTL,
		IsPublic:    request.IsPublic,
		PublicUsers: request.PublicUsers,
	}
}

func postRequestToNote(request PostRequest) notepkg.Note {
	return notepkg.Note{
		UserID:      request.UserID,
		Subject:     request.Subject,
		Text:        request.Text,
		TTL:         request.TTL,
		IsPublic:    request.IsPublic,
		PublicUsers: request.PublicUsers,
	}
}
