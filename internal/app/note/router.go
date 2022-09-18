package note

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"note-service/internal/app"
	notepkg "note-service/internal/pkg/note"
)

type noteService interface {
	CreateNote(note notepkg.Note) (notepkg.Note, error)
	FindNoteByID(id, userIDs string) (notepkg.Note, error)
	GetNotes(noteID string) ([]notepkg.Note, error)
	UpdateNote(note notepkg.Note) (notepkg.Note, error)
	DeleteNote(id, userID string) error
}
type Router struct {
	service noteService
}

func NewRouter(service noteService) *Router {
	return &Router{service}
}

func (r *Router) SetUpRouter(engine *gin.Engine) {
	engine.GET("/notes", app.AuthMiddleware(), r.getNotes)
	engine.GET("/note/:id", app.AuthMiddleware(), r.getNoteByID)
	engine.POST("/note", app.AuthMiddleware(), r.postNote)
	engine.PUT("/note/:id", app.AuthMiddleware(), r.updateNote)
	engine.DELETE("/note/:id", app.AuthMiddleware(), r.deleteNote)
}

func (r *Router) postNote(c *gin.Context) {
	var request PostRequest
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, app.ErrorModel{Error: err.Error()})
		return
	}

	request.UserID = c.GetString("userId")
	err := request.Validate()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	note := postRequestToNote(request)
	n, err := r.service.CreateNote(note)
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

	var request UpdateRequest
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, app.ErrorModel{Error: err.Error()})
		return
	}

	if request.ID != c.Param("id") {
		c.IndentedJSON(http.StatusBadRequest, app.ErrorModel{Error: "id in url and json not equal"})
		return
	}
	request.UserID = c.GetString("userId")
	err := request.Validate()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	note := updateRequestToNote(request)
	n, err := r.service.UpdateNote(note)
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
	c.IndentedJSON(http.StatusOK, n)
}

func (r *Router) deleteNote(c *gin.Context) {
	id := c.Param("id")
	UserID := c.GetString("userId")
	err := r.service.DeleteNote(id, UserID)
	if err != nil {
		if errors.Is(err, notepkg.ErrNoteNotFound) {
			c.IndentedJSON(http.StatusNotFound, app.ErrorModel{Error: err.Error()})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, app.ErrorModel{Error: err.Error()})
		}
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"note": "note successfully deleted"})
}

func (r *Router) getNotes(c *gin.Context) {
	userID := c.GetString("userId")
	notes, err := r.service.GetNotes(userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, app.UnknownError)
		return
	}
	c.IndentedJSON(http.StatusOK, notes)
}

func (r *Router) getNoteByID(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userId")
	m, err := r.service.FindNoteByID(id, userID)
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
