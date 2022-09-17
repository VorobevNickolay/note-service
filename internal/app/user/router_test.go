package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"note-service/internal/app"
	"note-service/internal/pkg/user"
)

type userStoreMock struct {
	CreateUserFunc                func(name, password string) (user.User, error)
	FindUserByIDFunc              func(id string) (user.User, error)
	FindUserByNameAndPasswordFunc func(name, password string) (user.User, error)
}

func (u *userStoreMock) CreateUser(_ context.Context, name, password string) (user.User, error) {
	return u.CreateUserFunc(name, password)
}

func (u *userStoreMock) FindUserByID(_ context.Context, id string) (user.User, error) {
	return u.FindUserByIDFunc(id)
}

func (u *userStoreMock) FindUserByNameAndPassword(_ context.Context, name, password string) (user.User, error) {
	return u.FindUserByNameAndPasswordFunc(name, password)
}

func TestGetUserById(t *testing.T) {
	tests := []struct {
		name              string
		userStore         userStoreMock
		userID            string
		expectedCode      int
		expectedUserModel Model
		expectedError     *app.ErrorModel
	}{
		{
			name: "should return errUserNotFound",
			userStore: userStoreMock{
				FindUserByIDFunc: func(id string) (user.User, error) {
					return user.User{}, user.ErrUserNotFound
				},
			},
			userID:        uuid.NewString(),
			expectedCode:  http.StatusNotFound,
			expectedError: &app.ErrorModel{Error: user.ErrUserNotFound.Error()},
		},
		{
			name: "should return user",
			userStore: userStoreMock{
				FindUserByIDFunc: func(id string) (user.User, error) {
					return user.User{ID: "ID1", Username: "User1", Password: "Password1"}, nil
				},
			},
			userID:            uuid.NewString(),
			expectedCode:      http.StatusOK,
			expectedUserModel: userModelFromUser(user.User{ID: "ID1", Username: "User1", Password: "Password1"}),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			r := NewRouter(&tt.userStore)
			r.SetUpRouter(g)

			req, _ := http.NewRequest(http.MethodGet, "/user/"+tt.userID, nil)
			w := httptest.NewRecorder()
			g.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			emptyUserModel := Model{}
			if tt.expectedUserModel != emptyUserModel {
				var actualUserModel Model
				err := json.Unmarshal(w.Body.Bytes(), &actualUserModel)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedUserModel, actualUserModel)
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

func TestSignUp(t *testing.T) {
	tests := []struct {
		name              string
		userStore         userStoreMock
		sentJSON          []byte
		expectedCode      int
		expectedUserModel Model
		expectedError     *app.ErrorModel
	}{
		{
			name: "should return DataBase error",
			userStore: userStoreMock{
				CreateUserFunc: func(name, password string) (user.User, error) {
					return user.User{}, errors.New("something wrong with db")
				},
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: &app.UnknownError,
		},
		{
			name: "should return ErrEmptyPassword",
			userStore: userStoreMock{
				CreateUserFunc: func(name, password string) (user.User, error) {
					return user.User{}, user.ErrEmptyPassword
				},
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: &app.ErrorModel{Error: user.ErrEmptyPassword.Error()},
		},
		{
			name: "should return ErrUsedUsername",
			userStore: userStoreMock{
				CreateUserFunc: func(name, password string) (user.User, error) {
					return user.User{}, user.ErrUsedUsername
				},
			},
			expectedCode:  http.StatusConflict,
			expectedError: &app.ErrorModel{Error: user.ErrUsedUsername.Error()},
		},
		{
			name: "should create user",
			userStore: userStoreMock{
				CreateUserFunc: func(name, password string) (user.User, error) {
					return user.User{ID: "ID1", Username: "Username1", Password: "Password1"}, nil
				},
			},
			expectedCode:      http.StatusCreated,
			expectedUserModel: userModelFromUser(user.User{ID: "ID1", Username: "Username1", Password: "Password1"}),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			r := NewRouter(&tt.userStore)
			r.SetUpRouter(g)

			var u = user.User{}

			jsonValue, _ := json.Marshal(u)
			req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			g.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			emptyUserModel := Model{}
			if tt.expectedUserModel != emptyUserModel {
				var actualUserModel Model
				err := json.Unmarshal(w.Body.Bytes(), &actualUserModel)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedUserModel, actualUserModel)
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
