package note

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"note-service/internal/app"
	"note-service/internal/pkg/jwt"
	"note-service/internal/pkg/note"
	"testing"
)

type noteServiceMock struct {
	CreateNoteFunc   func(note note.Note) (note.Note, error)
	FindNoteByIDFunc func(id, userIDs string) (note.Note, error)
	GetNotesFunc     func(userID, param string) ([]note.Note, error)
	UpdateNoteFunc   func(note note.Note) (note.Note, error)
	DeleteNoteFunc   func(id, userID string) error
}

func (n *noteServiceMock) CreateNote(note note.Note) (note.Note, error) {
	return n.CreateNoteFunc(note)
}

func (n *noteServiceMock) FindNoteByID(id, userID string) (note.Note, error) {
	return n.FindNoteByIDFunc(id, userID)
}

func (n *noteServiceMock) GetNotes(userID, param string) ([]note.Note, error) {
	return n.GetNotesFunc(userID, param)
}

func (n *noteServiceMock) UpdateNote(note note.Note) (note.Note, error) {
	return n.UpdateNoteFunc(note)
}

func (n *noteServiceMock) DeleteNote(id, userID string) error {
	return n.DeleteNoteFunc(id, userID)
}

func TestCreateNote(t *testing.T) {
	tests := []struct {
		name          string
		noteService   noteServiceMock
		Request       PostRequest
		expectedCode  int
		expectedError *app.ErrorModel
		expectedNote  NoteResponse
	}{
		{
			name: "should return request error",
			noteService: noteServiceMock{
				CreateNoteFunc: func(n note.Note) (note.Note, error) {
					return note.Note{}, nil
				},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "should return unknownError",
			Request: PostRequest{Text: "123", UserID: "123-123"},
			noteService: noteServiceMock{
				CreateNoteFunc: func(n note.Note) (note.Note, error) {
					return note.Note{}, errors.New("something wrong")
				},
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: &app.UnknownError,
		},
		{
			name:    "should return Note",
			Request: PostRequest{Text: "123", UserID: "123-123"},
			noteService: noteServiceMock{
				CreateNoteFunc: func(n note.Note) (note.Note, error) {
					return note.Note{ID: "123-123-123", Text: "123", UserID: "123-123"}, nil
				},
			},
			expectedCode: http.StatusCreated,
			expectedNote: noteToNoteResponse(note.Note{ID: "123-123-123", Text: "123", UserID: "123-123"}),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			logger, _ := zap.NewProduction()
			r := NewRouter(&tt.noteService, logger.Named(""))
			r.SetUpRouter(g)

			jsonValue, _ := json.Marshal(tt.Request)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequestWithContext(c, http.MethodPost, "/note", bytes.NewBuffer(jsonValue))
			token, _ := jwt.CreateToken("123-123")
			req.Header.Set(app.AccessHeader, token)
			g.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			emptyResponse := NoteResponse{}
			if tt.expectedNote != emptyResponse {
				var response NoteResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedNote, response)
			}
			if tt.expectedError != nil {
				var errorModel app.ErrorModel
				err := json.Unmarshal(w.Body.Bytes(), &errorModel)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedError, &errorModel)
			}
		})
	}
}

func TestUpdateNote(t *testing.T) {
	tests := []struct {
		name          string
		noteService   noteServiceMock
		Request       UpdateRequest
		id            string
		expectedCode  int
		expectedError *app.ErrorModel
		expectedNote  NoteResponse
	}{
		{
			name: "should return request error",
			id:   "123-123",
			noteService: noteServiceMock{
				UpdateNoteFunc: func(n note.Note) (note.Note, error) {
					return note.Note{}, nil
				},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "should return UrlIDError",
			Request: UpdateRequest{ID: "123-123", Text: "123"},
			id:      "123",
			noteService: noteServiceMock{
				UpdateNoteFunc: func(n note.Note) (note.Note, error) {
					return note.Note{}, errors.New("something wrong")
				},
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: &app.UrlIDError,
		},
		{
			name:    "should return errNoteNotFound",
			Request: UpdateRequest{ID: "123-123", Text: "123"},
			id:      "123-123",
			noteService: noteServiceMock{
				UpdateNoteFunc: func(n note.Note) (note.Note, error) {
					return note.Note{}, note.ErrNoteNotFound
				},
			},
			expectedCode:  http.StatusNotFound,
			expectedError: &app.ErrorModel{Error: note.ErrNoteNotFound.Error()},
		},
		{
			name:    "should return app.UnknownError",
			Request: UpdateRequest{ID: "123-123", Text: "123"},
			id:      "123-123",
			noteService: noteServiceMock{
				UpdateNoteFunc: func(n note.Note) (note.Note, error) {
					return note.Note{}, errors.New("something wrong")
				},
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: &app.UnknownError,
		},
		{
			name:    "should update Note",
			Request: UpdateRequest{ID: "123-123", Text: "123"},
			id:      "123-123",
			noteService: noteServiceMock{
				UpdateNoteFunc: func(n note.Note) (note.Note, error) {
					return note.Note{ID: "123-123", Text: "123"}, nil
				},
			},
			expectedCode: http.StatusOK,
			expectedNote: noteToNoteResponse(note.Note{ID: "123-123", Text: "123"}),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			logger, _ := zap.NewProduction()
			r := NewRouter(&tt.noteService, logger.Named(""))
			r.SetUpRouter(g)

			jsonValue, _ := json.Marshal(tt.Request)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequestWithContext(c, http.MethodPut, "/note/"+tt.id, bytes.NewBuffer(jsonValue))
			token, _ := jwt.CreateToken("123-123")
			req.Header.Set(app.AccessHeader, token)
			g.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			emptyResponse := NoteResponse{}
			if tt.expectedNote != emptyResponse {
				var response NoteResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedNote, response)
			}
			if tt.expectedError != nil {
				var errorModel app.ErrorModel
				err := json.Unmarshal(w.Body.Bytes(), &errorModel)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedError, &errorModel)
			}
		})
	}
}

func TestDeleteNote(t *testing.T) {
	tests := []struct {
		name          string
		noteService   noteServiceMock
		id            string
		expectedCode  int
		expectedError *app.ErrorModel
		expectedNote  NoteResponse
	}{
		{
			name: "should return errorNotFound",
			id:   "123-123",
			noteService: noteServiceMock{
				DeleteNoteFunc: func(id, userID string) error {
					return note.ErrNoteNotFound
				},
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "should return unknownError",
			id:   "123-123",
			noteService: noteServiceMock{
				DeleteNoteFunc: func(id, userID string) error {
					return errors.New("something wrong")
				},
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: &app.UnknownError,
		},
		{
			name: "should delete message",
			id:   "123-123",
			noteService: noteServiceMock{
				DeleteNoteFunc: func(id, userID string) error {
					return nil
				},
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			logger, _ := zap.NewProduction()
			r := NewRouter(&tt.noteService, logger.Named(""))
			r.SetUpRouter(g)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequestWithContext(c, http.MethodDelete, "/note/"+tt.id, nil)
			token, _ := jwt.CreateToken("123-123")
			req.Header.Set(app.AccessHeader, token)
			g.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			emptyResponse := NoteResponse{}
			if tt.expectedNote != emptyResponse {
				var response NoteResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedNote, response)
			}
			if tt.expectedError != nil {
				var errorModel app.ErrorModel
				err := json.Unmarshal(w.Body.Bytes(), &errorModel)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedError, &errorModel)
			}
		})
	}
}

func TestGetNotes(t *testing.T) {
	tests := []struct {
		name          string
		noteService   noteServiceMock
		param         string
		expectedCode  int
		expectedError *app.ErrorModel
		expectedNote  []note.Note
	}{
		{
			name: "should return unknownError",
			noteService: noteServiceMock{
				GetNotesFunc: func(userID, param string) ([]note.Note, error) {
					return []note.Note{}, errors.New("something wrong")
				},
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: &app.UnknownError,
		},
		{
			name: "should return Notes",
			noteService: noteServiceMock{
				GetNotesFunc: func(userID, param string) ([]note.Note, error) {
					return []note.Note{{ID: "123-123", Text: "123-123"}, {ID: "123-124", Text: "123-123"}}, nil
				},
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			logger, _ := zap.NewProduction()
			r := NewRouter(&tt.noteService, logger.Named(""))
			r.SetUpRouter(g)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			jsonValue, _ := json.Marshal(tt.param)
			req, _ := http.NewRequestWithContext(c, http.MethodGet, "/notes", bytes.NewBuffer(jsonValue))
			token, _ := jwt.CreateToken("123-123")
			req.Header.Set(app.AccessHeader, token)
			g.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if len(tt.expectedNote) != 0 {
				var response NoteResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedNote, response)
			}
			if tt.expectedError != nil {
				var errorModel app.ErrorModel
				err := json.Unmarshal(w.Body.Bytes(), &errorModel)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedError, &errorModel)
			}
		})
	}
}

func TestGetNoteByID(t *testing.T) {
	tests := []struct {
		name          string
		noteService   noteServiceMock
		id            string
		expectedCode  int
		expectedError *app.ErrorModel
		expectedNote  NoteResponse
	}{
		{
			name: "should return errNoteNotFound",
			id:   "123-123",
			noteService: noteServiceMock{
				FindNoteByIDFunc: func(id, userID string) (note.Note, error) {
					return note.Note{}, note.ErrNoteNotFound
				},
			},
			expectedCode:  http.StatusNotFound,
			expectedError: &app.ErrorModel{Error: note.ErrNoteNotFound.Error()},
		},
		{
			name: "should return unknownError",
			id:   "123-123",
			noteService: noteServiceMock{
				FindNoteByIDFunc: func(id, userID string) (note.Note, error) {
					return note.Note{}, errors.New("something wrong")
				},
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: &app.UnknownError,
		},
		{
			name: "should return Note",
			id:   "123-123",
			noteService: noteServiceMock{
				FindNoteByIDFunc: func(id, userID string) (note.Note, error) {
					return note.Note{ID: "123-123", Text: "123"}, nil
				},
			},
			expectedCode: http.StatusOK,
			expectedNote: noteToNoteResponse(note.Note{ID: "123-123", Text: "123"}),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			logger, _ := zap.NewProduction()
			r := NewRouter(&tt.noteService, logger.Named(""))
			r.SetUpRouter(g)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequestWithContext(c, http.MethodGet, "/note/"+tt.id, nil)
			token, _ := jwt.CreateToken("123-123")
			req.Header.Set(app.AccessHeader, token)
			g.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			emptyResponse := NoteResponse{}
			if tt.expectedNote != emptyResponse {
				var response NoteResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedNote, response)
			}
			if tt.expectedError != nil {
				var errorModel app.ErrorModel
				err := json.Unmarshal(w.Body.Bytes(), &errorModel)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedError, &errorModel)
			}
		})
	}
}
