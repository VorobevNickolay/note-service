package note

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"note-service/internal/app"
	notepkg "note-service/internal/pkg/note"
)

type noteStore interface {
	CreateNote(ctx context.Context, note notepkg.Note) (notepkg.Note, error)
	FindNoteByID(ctx context.Context, id string) (notepkg.Note, error)
	GetNotes(ctx context.Context, noteID string) ([]notepkg.Note, error)
	UpdateNote(ctx context.Context, note notepkg.Note) (notepkg.Note, error)
	DeleteNote(ctx context.Context, id string) error
}
type Router struct {
	store noteStore
}

func NewRouter(store noteStore) *Router {
	return &Router{store}
}

func (r *Router) SetUpRouter(engine *gin.Engine) {
	engine.GET("/notes", app.AuthMiddleware(), r.getNotes)
	engine.GET("/note/:id", app.AuthMiddleware(), r.getNoteByID)
	engine.POST("/note", app.AuthMiddleware(), r.postNote)
	engine.PUT("/note/:id", app.AuthMiddleware(), r.updateNote)
	engine.DELETE("/note/:id", app.AuthMiddleware(), r.deleteNote)
}

func (r *Router) postNote(c *gin.Context) {
	var newNote notepkg.Note
	if err := c.BindJSON(&newNote); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, app.ErrorModel{Error: err.Error()})
		return
	}

	newNote.UserID = c.GetString("userId")
	n, err := r.store.CreateNote(c, newNote)
	if err != nil {
		if errors.Is(err, notepkg.ErrEmptyNote) {
			c.IndentedJSON(http.StatusBadRequest, app.ErrorModel{Error: err.Error()})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, app.UnknownError)
		return
	}
	c.IndentedJSON(http.StatusCreated, n)
}

func (r *Router) updateNote(c *gin.Context) {
	id := c.Param("id")
	noteID := c.GetString("noteID")
	oldNote, err := r.store.FindNoteByID(c, id)
	if err != nil {
		if errors.Is(err, notepkg.ErrNoteNotFound) {
			c.IndentedJSON(http.StatusNotFound, app.ErrorModel{Error: err.Error()})
		}
		c.IndentedJSON(http.StatusInternalServerError, app.UnknownError)
		return
	}

	if oldNote.ID != noteID {
		c.IndentedJSON(http.StatusForbidden, app.ErrorModel{Error: app.ErrNoAccess.Error()})
		return
	}

	var newNote notepkg.Note
	if err := c.BindJSON(&newNote); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, app.ErrorModel{Error: err.Error()})
		return
	}

	m, err := r.store.UpdateNote(c, newNote)
	if err != nil {
		if errors.Is(err, notepkg.ErrNoteNotFound) {
			c.IndentedJSON(http.StatusNotFound, app.ErrorModel{Error: err.Error()})
		} else if errors.Is(err, notepkg.ErrEmptyNote) {
			c.IndentedJSON(http.StatusBadRequest, app.ErrorModel{Error: err.Error()})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, app.UnknownError)
		}
		return
	}
	c.IndentedJSON(http.StatusOK, m)
}

func (r *Router) deleteNote(c *gin.Context) {
	id := c.Param("id")
	noteID := c.GetString("noteID")
	oldNote, err := r.store.FindNoteByID(c, id)
	if err != nil {
		if errors.Is(err, notepkg.ErrNoteNotFound) {
			c.IndentedJSON(http.StatusNotFound, app.ErrorModel{Error: err.Error()})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, app.UnknownError)
		}
		return
	}

	if oldNote.ID != noteID {
		c.IndentedJSON(http.StatusForbidden, app.ErrorModel{Error: app.ErrNoAccess.Error()})
		return
	}

	err = r.store.DeleteNote(c, id)
	if err != nil {
		if errors.Is(err, notepkg.ErrNoteNotFound) {
			c.IndentedJSON(http.StatusNotFound, app.ErrorModel{Error: err.Error()})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, app.UnknownError)
		}
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"note": "note successfully deleted"})
}

func (r *Router) getNotes(c *gin.Context) {
	userID := c.GetString("userId")
	notes, err := r.store.GetNotes(c, userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, app.UnknownError)
		return
	}
	c.IndentedJSON(http.StatusOK, notes)
}

func (r *Router) getNoteByID(c *gin.Context) {
	id := c.Param("id")

	m, err := r.store.FindNoteByID(c, id)
	if err != nil {
		if errors.Is(err, notepkg.ErrNoteNotFound) {
			c.IndentedJSON(http.StatusNotFound, app.ErrorModel{Error: err.Error()})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, app.UnknownError)
		}
		return
	}

	c.IndentedJSON(http.StatusOK, m)
}
