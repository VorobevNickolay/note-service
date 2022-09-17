package note

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/exp/maps"
	"sync"
	"time"
)

// notes map[userId]map[noteId]Note
// noteIDs map[noteId] userId
type inMemoryStore struct {
	sync.RWMutex
	notes   map[string]map[string]Note
	noteIDs map[string]string
}

func NewInMemoryStore() *inMemoryStore {
	return &inMemoryStore{
		notes:   make(map[string]map[string]Note, 0),
		noteIDs: make(map[string]string, 0),
	}
}

func (store *inMemoryStore) CreateNote(_ context.Context, note Note) (Note, error) {
	store.Lock()
	defer store.Unlock()

	note.ID = uuid.NewString()
	note.CreatedAt = time.Now().UTC()
	if note.Text == "" {
		return Note{}, ErrEmptyNote
	}
	if _, ok := store.notes[note.UserID]; !ok {
		store.notes[note.UserID] = make(map[string]Note, 0)
	}
	store.notes[note.UserID][note.ID] = note
	store.noteIDs[note.ID] = note.UserID
	return note, nil
}

func (store *inMemoryStore) GetNotes(_ context.Context, userID string) ([]Note, error) {
	v := maps.Values(store.notes[userID])
	return v, nil
}
func (store *inMemoryStore) FindNoteByID(_ context.Context, id string) (Note, error) {
	store.RLock()
	defer store.RUnlock()

	if userID, ok := store.noteIDs[id]; ok {
		return store.notes[userID][id], nil
	}
	return Note{}, ErrNoteNotFound
}

func (store *inMemoryStore) DeleteNote(_ context.Context, id string) error {
	store.Lock()
	defer store.Unlock()

	userID, ok := store.noteIDs[id]
	if !ok {
		return ErrNoteNotFound
	}
	delete(store.notes[userID], id)
	delete(store.noteIDs, id)
	return nil
}

func (store *inMemoryStore) UpdateNote(_ context.Context, note Note) (Note, error) {
	store.Lock()
	defer store.Unlock()
	var oldNote *Note

	userID, ok := store.noteIDs[note.ID]
	if !ok {
		return Note{}, ErrNoteNotFound
	}
	*oldNote = store.notes[userID][note.ID]
	*oldNote = note
	oldNote.UpdatedAt = time.Now().UTC()

	return *oldNote, nil
}
