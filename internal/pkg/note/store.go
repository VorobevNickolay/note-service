package note

import (
	"go.uber.org/zap"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"
)

// notes map[userId]map[noteId]Note
// noteIDs map[noteId] userId
type InMemoryStore struct {
	sync.RWMutex
	notes    map[string]map[string]Note
	noteIDs  map[string]string
	notesTTL map[string]int64
	logger   *zap.Logger
}

func NewInMemoryStore(logger *zap.Logger) *InMemoryStore {
	return &InMemoryStore{
		notes:    make(map[string]map[string]Note, 0),
		noteIDs:  make(map[string]string, 0),
		notesTTL: make(map[string]int64, 0),
		logger:   logger,
	}
}

func (store *InMemoryStore) CreateNote(note Note) (Note, error) {
	store.Lock()
	defer store.Unlock()

	note.ID = uuid.NewString()
	note.CreatedAt = time.Now().UTC()
	if _, ok := store.notes[note.UserID]; !ok {
		store.notes[note.UserID] = make(map[string]Note, 0)
	}
	store.notes[note.UserID][note.ID] = note
	store.noteIDs[note.ID] = note.UserID
	if note.TTL != nil {
		store.notesTTL[note.ID] = *note.TTL
	}

	return note, nil
}

func (store *InMemoryStore) GetNotes(userID, param string) ([]Note, error) {
	v := maps.Values(store.notes[userID])
	switch param {
	case "ttl":
		sort.SliceStable(v, func(i, j int) bool {
			return *v[i].TTL < *v[j].TTL
		})
	case "subject":
		sort.SliceStable(v, func(i, j int) bool {
			return strings.Compare(v[i].Subject, v[j].Subject) < 0
		})
	case "created-at":
		sort.SliceStable(v, func(i, j int) bool {
			return v[i].CreatedAt.Before(v[j].CreatedAt)
		})
	case "updated-at":
		sort.SliceStable(v, func(i, j int) bool {
			return v[i].UpdatedAt.Before(v[j].UpdatedAt)
		})
	default:
	}

	return v, nil
}

func (store *InMemoryStore) FindNoteByID(id string) (Note, error) {
	store.RLock()
	defer store.RUnlock()

	if userID, ok := store.noteIDs[id]; ok {
		return store.notes[userID][id], nil
	}
	return Note{}, ErrNoteNotFound
}

func (store *InMemoryStore) DeleteNote(id string) error {
	store.Lock()
	defer store.Unlock()

	userID, ok := store.noteIDs[id]
	if !ok {
		return ErrNoteNotFound
	}
	delete(store.notes[userID], id)
	delete(store.noteIDs, id)
	delete(store.notesTTL, id)
	return nil
}

func (store *InMemoryStore) UpdateNote(note Note) (Note, error) {
	store.Lock()
	defer store.Unlock()

	note.UpdatedAt = time.Now().UTC()
	store.notes[note.UserID][note.ID] = note
	if note.TTL != nil {
		store.notesTTL[note.ID] = *note.TTL
	} else {
		delete(store.notesTTL, note.ID)
	}

	return note, nil
}

func (store *InMemoryStore) ExpireNotes() error {
	store.Lock()
	defer store.Unlock()

	for noteID, ttl := range store.notesTTL {
		if time.Now().UTC().Unix() >= ttl {
			userID := store.noteIDs[noteID]
			delete(store.notes[userID], noteID)
			delete(store.noteIDs, noteID)
			delete(store.notesTTL, noteID)
			store.logger.Info("note was deleted", zap.String("noteID", noteID))
		}
	}

	return nil
}
