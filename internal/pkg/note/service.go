package note

import (
	"note-service/internal/app"
)

type store interface {
	CreateNote(note Note) (Note, error)
	FindNoteByID(id string) (Note, error)
	GetNotes(userID, param string) ([]Note, error)
	UpdateNote(note Note) (Note, error)
	DeleteNote(id string) error
}

type Service struct {
	store store
}

func NewService(store store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateNote(note Note) (Note, error) {
	return s.store.CreateNote(note)
}

func (s *Service) FindNoteByID(id, userID string) (Note, error) {
	note, err := s.store.FindNoteByID(id)
	if err != nil {
		return Note{}, err
	}

	if userID == note.UserID || (note.IsPublic && note.PublicUsers == nil) {
		return note, nil
	}

	pu := note.PublicUsers
	for _, u := range *pu {
		if u == userID {
			return note, nil
		}
	}
	return Note{}, app.ErrNoAccess
}

func (s *Service) GetNotes(id, param string) ([]Note, error) {
	return s.store.GetNotes(id, param)
}

func (s *Service) UpdateNote(note Note) (Note, error) {
	n, err := s.store.FindNoteByID(note.ID)
	if err != nil {
		return Note{}, err
	}

	if n.UserID != note.UserID {
		return Note{}, app.ErrNoAccess
	}
	note.CreatedAt = n.CreatedAt
	return s.store.UpdateNote(note)
}

func (s *Service) DeleteNote(id, userID string) error {
	n, err := s.store.FindNoteByID(id)
	if err != nil {
		return err
	}
	if n.UserID != userID {
		return app.ErrNoAccess
	}
	return s.store.DeleteNote(id)
}
