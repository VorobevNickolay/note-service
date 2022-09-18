package note

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"
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

func (store *inMemoryStore) CreateNote(note Note) (Note, error) {
	store.Lock()
	defer store.Unlock()

	note.ID = uuid.NewString()
	note.CreatedAt = time.Now().UTC()
	if _, ok := store.notes[note.UserID]; !ok {
		store.notes[note.UserID] = make(map[string]Note, 0)
	}
	store.notes[note.UserID][note.ID] = note
	store.noteIDs[note.ID] = note.UserID
	return note, nil
}

func (store *inMemoryStore) GetNotes(userID string) ([]Note, error) {
	v := maps.Values(store.notes[userID])
	return v, nil
}

func (store *inMemoryStore) FindNoteByID(id string) (Note, error) {
	store.RLock()
	defer store.RUnlock()

	if userID, ok := store.noteIDs[id]; ok {
		return store.notes[userID][id], nil
	}
	return Note{}, ErrNoteNotFound
}

func (store *inMemoryStore) DeleteNote(id string) error {
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

func (store *inMemoryStore) UpdateNote(note Note) (Note, error) {
	store.Lock()
	defer store.Unlock()

	note.UpdatedAt = time.Now().UTC()
	store.notes[note.UserID][note.ID] = note

	return note, nil
}
