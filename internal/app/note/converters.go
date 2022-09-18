package note

import (
	notepkg "note-service/internal/pkg/note"
)

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

func noteToNoteResponse(note notepkg.Note) NoteResponse {
	return NoteResponse{
		ID:          note.ID,
		UserID:      note.UserID,
		Subject:     note.Subject,
		Text:        note.Text,
		TTL:         note.TTL,
		IsPublic:    note.IsPublic,
		PublicUsers: note.PublicUsers,
		CreatedAt:   note.CreatedAt,
		UpdatedAt:   note.UpdatedAt,
	}
}

func notesToNoteResponses(notes []notepkg.Note) []NoteResponse {
	res := make([]NoteResponse, len(notes))
	for i, note := range notes {
		res[i] = noteToNoteResponse(note)
	}
	return res
}
