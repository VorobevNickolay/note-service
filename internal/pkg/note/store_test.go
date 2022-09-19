package note

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestGetNotes(t *testing.T) {
	t.Run("should return empty list", func(t *testing.T) {
		store := NewInMemoryStore(&zap.Logger{})
		actual, err := store.GetNotes("123-123-123", "")
		require.NoError(t, err)
		require.Equal(t, 0, len(actual))
	})

	t.Run("should return notes sort by createdAt", func(t *testing.T) {
		store := NewInMemoryStore(&zap.Logger{})
		note1 := Note{Text: "123-123", UserID: "123-123-123"}
		note2 := Note{Text: "123-123", UserID: "123-123-123"}
		note1, err := store.CreateNote(note1)
		require.NoError(t, err)

		note2, err = store.CreateNote(note2)
		require.NoError(t, err)

		actual, err := store.GetNotes("123-123-123", "created-at")
		require.NoError(t, err)

		require.Equal(t, 2, len(actual))
		require.Equal(t, note1, actual[0])
		require.Equal(t, note2, actual[1])
	})
	t.Run("should return notes sort by Subject", func(t *testing.T) {
		store := NewInMemoryStore(&zap.Logger{})
		note1 := Note{Text: "123-123", UserID: "123-123-123", Subject: "Ca"}
		note2 := Note{Text: "123-123", UserID: "123-123-123", Subject: "Ab"}
		note3 := Note{Text: "123-123", UserID: "123-123-123", Subject: "Ea"}
		note1, err := store.CreateNote(note1)
		require.NoError(t, err)

		note2, err = store.CreateNote(note2)
		require.NoError(t, err)

		note3, err = store.CreateNote(note3)
		require.NoError(t, err)

		actual, err := store.GetNotes("123-123-123", "subject")
		require.NoError(t, err)

		require.Equal(t, 3, len(actual))
		require.Equal(t, note1, actual[1])
		require.Equal(t, note2, actual[0])
		require.Equal(t, note3, actual[2])
	})
	t.Run("should return notes sort by ttl", func(t *testing.T) {
		store := NewInMemoryStore(&zap.Logger{})
		ttl1, ttl2, ttl3 := int64(10), int64(20), int64(30)
		note1 := Note{Text: "123-123", UserID: "123-123-123", TTL: &ttl3}
		note2 := Note{Text: "123-123", UserID: "123-123-123", TTL: &ttl2}
		note3 := Note{Text: "123-123", UserID: "123-123-123", TTL: &ttl1}
		note1, err := store.CreateNote(note1)
		require.NoError(t, err)

		note2, err = store.CreateNote(note2)
		require.NoError(t, err)

		note3, err = store.CreateNote(note3)
		require.NoError(t, err)

		actual, err := store.GetNotes("123-123-123", "ttl")
		require.NoError(t, err)

		require.Equal(t, 3, len(actual))
		require.Equal(t, note1, actual[2])
		require.Equal(t, note2, actual[1])
		require.Equal(t, note3, actual[0])
	})
	t.Run("should return notes sort by updated-at", func(t *testing.T) {
		store := NewInMemoryStore(&zap.Logger{})
		note1 := Note{Text: "123-123", UserID: "123-123-123"}
		note2 := Note{Text: "123-123", UserID: "123-123-123"}
		note3 := Note{Text: "123-123", UserID: "123-123-123"}
		note1, err := store.CreateNote(note1)
		require.NoError(t, err)

		note2, err = store.CreateNote(note2)
		require.NoError(t, err)

		note3, err = store.CreateNote(note3)
		require.NoError(t, err)

		ttl := int64(10)
		note3.TTL = &ttl
		note3, err = store.UpdateNote(note3)
		require.NoError(t, err)

		note1, err = store.UpdateNote(note1)
		require.NoError(t, err)

		note2, err = store.UpdateNote(note2)
		require.NoError(t, err)

		actual, err := store.GetNotes("123-123-123", "updated-at")
		require.NoError(t, err)

		require.Equal(t, 3, len(actual))
		require.Equal(t, note1, actual[1])
		require.Equal(t, note2, actual[2])
		require.Equal(t, note3, actual[0])
	})
}
func TestFindNoteID(t *testing.T) {
	t.Run("should return errNoteNotFound", func(t *testing.T) {
		store := NewInMemoryStore(&zap.Logger{})
		actual, err := store.FindNoteByID("123-123-123")
		require.Error(t, err, ErrNoteNotFound)
		require.Empty(t, actual)
	})

	t.Run("should return note", func(t *testing.T) {
		store := NewInMemoryStore(&zap.Logger{})
		note1 := Note{Text: "123-123", UserID: "123-123-123"}
		note2 := Note{Text: "123-123", UserID: "123-123-123"}
		note1, err := store.CreateNote(note1)
		require.NoError(t, err)

		note2, err = store.CreateNote(note2)
		require.NoError(t, err)

		actual, err := store.FindNoteByID(note1.ID)
		require.NoError(t, err)

		require.Equal(t, note1, actual)
	})
}

func TestDeleteNote(t *testing.T) {
	t.Run("should return errNoteNotFound", func(t *testing.T) {
		store := NewInMemoryStore(&zap.Logger{})
		err := store.DeleteNote("123-123-123")
		require.Error(t, err, ErrNoteNotFound)
	})

	t.Run("should delete note", func(t *testing.T) {
		store := NewInMemoryStore(&zap.Logger{})
		note1 := Note{Text: "123-123", UserID: "123-123-123"}
		note2 := Note{Text: "123-123", UserID: "123-123-123"}
		note1, err := store.CreateNote(note1)
		require.NoError(t, err)

		note2, err = store.CreateNote(note2)
		require.NoError(t, err)

		err = store.DeleteNote(note1.ID)
		require.NoError(t, err)
		actual, _ := store.GetNotes("123-123-123", "")
		require.Equal(t, note2, actual[0])
	})
}

func TestExpireNotes(t *testing.T) {
	t.Run("should delete note", func(t *testing.T) {
		logger, _ := zap.NewProduction()
		store := NewInMemoryStore(logger.Named("test-logger"))
		ttl := time.Now().UTC().Unix() + 5
		note1 := Note{Text: "123-123", UserID: "123-123-123", TTL: &ttl}
		note2 := Note{Text: "123-123", UserID: "123-123-123"}
		note1, err := store.CreateNote(note1)
		require.NoError(t, err)

		note2, err = store.CreateNote(note2)
		require.NoError(t, err)

		time.Sleep(7 * time.Second)
		err = store.ExpireNotes()
		require.NoError(t, err)
		actual, _ := store.GetNotes("123-123-123", "")
		require.Equal(t, note2, actual[0])
	})
}
