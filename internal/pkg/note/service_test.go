package note

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"note-service/internal/app"
	"testing"
)

type noteStoreMock struct {
	CreateNoteFunc   func(note Note) (Note, error)
	FindNoteByIDFunc func(id string) (Note, error)
	GetNotesFunc     func(userID, param string) ([]Note, error)
	UpdateNoteFunc   func(note Note) (Note, error)
	DeleteNoteFunc   func(id string) error
}

func (s *noteStoreMock) CreateNote(note Note) (Note, error) {
	return s.CreateNoteFunc(note)
}

func (s *noteStoreMock) FindNoteByID(id string) (Note, error) {
	return s.FindNoteByIDFunc(id)
}

func (s *noteStoreMock) GetNotes(userID, param string) ([]Note, error) {
	return s.GetNotesFunc(userID, param)
}

func (s *noteStoreMock) UpdateNote(note Note) (Note, error) {
	return s.UpdateNoteFunc(note)
}

func (s *noteStoreMock) DeleteNote(id string) error {
	return s.DeleteNoteFunc(id)
}

func TestServiceGetNotes(t *testing.T) {
	tests := []struct {
		name          string
		noteStore     noteStoreMock
		id            string
		param         string
		expectedNotes []Note
		expectedError error
	}{
		{
			name: "should return notes",
			noteStore: noteStoreMock{
				GetNotesFunc: func(userID, param string) ([]Note, error) {
					return []Note{{ID: "123", Text: "123"}}, nil
				},
			},
			param:         uuid.NewString(),
			id:            uuid.NewString(),
			expectedNotes: []Note{{ID: "123", Text: "123"}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(&tt.noteStore)
			n, err := s.GetNotes(tt.id, tt.param)
			if len(tt.expectedNotes) != 0 {
				require.Equal(t, n, tt.expectedNotes)
			}
			if tt.expectedError != nil {
				require.Error(t, err, tt.expectedError)
			}

		})
	}
}

func TestServiceCreateNote(t *testing.T) {
	tests := []struct {
		name          string
		noteStore     noteStoreMock
		note          Note
		expectedNote  Note
		expectedError error
	}{
		{
			name: "should return notes",
			noteStore: noteStoreMock{
				CreateNoteFunc: func(note Note) (Note, error) {
					return Note{ID: "123", Text: "123"}, nil
				},
			},
			note:         Note{ID: "123", Text: "123"},
			expectedNote: Note{ID: "123", Text: "123"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(&tt.noteStore)
			n, err := s.CreateNote(tt.note)
			require.Equal(t, tt.expectedNote, n)
			if tt.expectedError != nil {
				require.Error(t, err, tt.expectedError)
			}

		})
	}
}

func TestFindNoteByID(t *testing.T) {
	tests := []struct {
		name          string
		noteStore     noteStoreMock
		id            string
		userID        string
		expectedNote  Note
		expectedError error
	}{
		{
			name:   "should return errNoteNotFound",
			id:     uuid.NewString(),
			userID: uuid.NewString(),
			noteStore: noteStoreMock{
				FindNoteByIDFunc: func(id string) (Note, error) {
					return Note{}, ErrNoteNotFound
				},
			},
			expectedError: ErrNoteNotFound,
		},
		{
			name:   "should return Note,public",
			id:     uuid.NewString(),
			userID: uuid.NewString(),
			noteStore: noteStoreMock{
				FindNoteByIDFunc: func(id string) (Note, error) {
					return Note{ID: "123-123-123", Text: "123", IsPublic: true}, nil
				},
			},
			expectedNote: Note{ID: "123-123-123", Text: "123", IsPublic: true},
		},
		{
			name:   "should return note, for invited user",
			id:     uuid.NewString(),
			userID: "123-123-123",
			noteStore: noteStoreMock{
				FindNoteByIDFunc: func(id string) (Note, error) {
					return Note{
						ID:          "123-123-123",
						Text:        "123",
						IsPublic:    true,
						PublicUsers: &[]string{"123-321-123", "123-123-123"}}, nil
				},
			},
			expectedNote: Note{
				ID:          "123-123-123",
				Text:        "123",
				IsPublic:    true,
				PublicUsers: &[]string{"123-321-123", "123-123-123"}},
		},
		{
			name:   "should return ErrNoAccess",
			id:     uuid.NewString(),
			userID: uuid.NewString(),
			noteStore: noteStoreMock{
				FindNoteByIDFunc: func(id string) (Note, error) {
					return Note{
						ID:          "123-123-123",
						Text:        "123",
						IsPublic:    true,
						PublicUsers: &[]string{"123-321-123", "123-123-123"}}, nil
				},
			},
			expectedError: app.ErrNoAccess,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(&tt.noteStore)
			n, err := s.FindNoteByID(tt.id, tt.userID)
			require.Equal(t, tt.expectedNote, n)
			if tt.expectedError != nil {
				require.Error(t, err, tt.expectedError)
			}

		})
	}
}

func TestServiceDeleteNote(t *testing.T) {
	tests := []struct {
		name          string
		noteStore     noteStoreMock
		id            string
		userID        string
		expectedError error
	}{
		{
			name: "should return errNoteNotFound",
			id:   uuid.NewString(),
			noteStore: noteStoreMock{
				FindNoteByIDFunc: func(id string) (Note, error) {
					return Note{}, ErrNoteNotFound
				},
			},
			expectedError: ErrNoteNotFound,
		},
		{
			name:   "should return errNoAccess",
			id:     uuid.NewString(),
			userID: uuid.NewString(),
			noteStore: noteStoreMock{
				FindNoteByIDFunc: func(id string) (Note, error) {
					return Note{UserID: "User2"}, nil
				},
			},
			expectedError: app.ErrNoAccess,
		},
		{
			name:   "should delete Note",
			id:     uuid.NewString(),
			userID: "User1",
			noteStore: noteStoreMock{
				FindNoteByIDFunc: func(id string) (Note, error) {
					return Note{UserID: "User1", Text: "123"}, nil
				},
				DeleteNoteFunc: func(id string) error {
					return nil
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(&tt.noteStore)
			err := s.DeleteNote(tt.id, tt.userID)
			if tt.expectedError != nil {
				require.Error(t, err, tt.expectedError)
			}

		})
	}
}

func TestUpdateNote(t *testing.T) {
	tests := []struct {
		name          string
		noteStore     noteStoreMock
		note          Note
		expectedNote  Note
		expectedError error
	}{
		{
			name: "should return errNoteNotFound",
			note: Note{ID: "123-123", Text: "123"},
			noteStore: noteStoreMock{
				FindNoteByIDFunc: func(id string) (Note, error) {
					return Note{}, ErrNoteNotFound
				},
			},
			expectedError: ErrNoteNotFound,
		},
		{
			name: "should return errNoAccess",
			note: Note{ID: "123-123", Text: "123", UserID: "User1"},
			noteStore: noteStoreMock{
				FindNoteByIDFunc: func(id string) (Note, error) {
					return Note{UserID: "User2"}, nil
				},
			},
			expectedError: app.ErrNoAccess,
		},
		{
			name: "should update Note",
			note: Note{ID: "123-123", Text: "123", UserID: "User2"},
			noteStore: noteStoreMock{
				FindNoteByIDFunc: func(id string) (Note, error) {
					return Note{UserID: "User2"}, nil
				},
				UpdateNoteFunc: func(note Note) (Note, error) {
					return Note{ID: "123-123-123", UserID: "User2", Text: "123"}, nil
				},
			},
			expectedNote: Note{ID: "123-123-123", UserID: "User2", Text: "123"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(&tt.noteStore)
			n, err := s.UpdateNote(tt.note)
			require.Equal(t, tt.expectedNote, n)
			if tt.expectedError != nil {
				require.Error(t, err, tt.expectedError)
			}

		})
	}
}
