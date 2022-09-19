package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

type userStoreMock struct {
	CreateUserFunc     func(name, password string) (User, error)
	FindUserByNameFunc func(name string) (User, error)
}

func (s *userStoreMock) CreateUser(name, password string) (User, error) {
	return s.CreateUserFunc(name, password)
}

func (s *userStoreMock) FindUserByName(name string) (User, error) {
	return s.FindUserByNameFunc(name)
}

func TestSignUp(t *testing.T) {
	tests := []struct {
		name          string
		userStore     userStoreMock
		username      string
		password      string
		expectedUser  User
		expectedError error
	}{
		{
			name: "should return errUsedUsername",
			userStore: userStoreMock{
				FindUserByNameFunc: func(id string) (User, error) {
					return User{ID: "123", Password: "123"}, nil
				},
			},
			username:      uuid.NewString(),
			password:      uuid.NewString(),
			expectedError: ErrUsedUsername,
		},
		{
			name: "should return createUserError",
			userStore: userStoreMock{
				FindUserByNameFunc: func(id string) (User, error) {
					return User{}, ErrUserNotFound
				},
				CreateUserFunc: func(name, password string) (User, error) {
					return User{}, errors.New("createUserError")
				},
			},
			username:      uuid.NewString(),
			password:      uuid.NewString(),
			expectedError: errors.New("createUserError"),
		},
		{
			name: "should createUser",
			userStore: userStoreMock{
				FindUserByNameFunc: func(id string) (User, error) {
					return User{}, ErrUserNotFound
				},
				CreateUserFunc: func(name, password string) (User, error) {
					return User{ID: "123-123-123", Username: "username1", Password: "password1"}, nil
				},
			},
			username:     uuid.NewString(),
			password:     uuid.NewString(),
			expectedUser: User{ID: "123-123-123", Username: "username1", Password: "password1"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(&tt.userStore)
			u, err := s.SignUp(tt.username, tt.password)
			emptyUser := User{}
			if tt.expectedUser != emptyUser {
				require.Equal(t, u, tt.expectedUser)
			}
			if tt.expectedError != nil {
				require.Error(t, err, tt.expectedError)
			}

		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name          string
		userStore     userStoreMock
		username      string
		password      string
		expectedUser  User
		expectedError error
	}{
		{
			name: "should return errUserNotFound from FindUserByName",
			userStore: userStoreMock{
				FindUserByNameFunc: func(id string) (User, error) {
					return User{}, ErrUserNotFound
				},
			},
			username:      uuid.NewString(),
			password:      uuid.NewString(),
			expectedError: ErrUserNotFound,
		},
		{
			name: "should return errUserNotFound, wrong password",
			userStore: userStoreMock{
				FindUserByNameFunc: func(id string) (User, error) {
					return User{"123-123", "username1", "password1"}, nil
				},
			},
			username:      uuid.NewString(),
			password:      uuid.NewString(),
			expectedError: ErrUserNotFound,
		},
		{
			name: "should login user",
			userStore: userStoreMock{
				FindUserByNameFunc: func(id string) (User, error) {
					return User{ID: "123-123-123", Username: "username1", Password: "$2a$10$7Fy455pjoxYl4f3.TGiPNut/pHy/K0C93oSwqkX.pDEDxGNvplrUG"}, nil
				},
			},
			username:     uuid.NewString(),
			password:     "123123123",
			expectedUser: User{ID: "123-123-123", Username: "username1", Password: "$2a$10$7Fy455pjoxYl4f3.TGiPNut/pHy/K0C93oSwqkX.pDEDxGNvplrUG"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(&tt.userStore)
			u, err := s.Login(tt.username, tt.password)
			emptyUser := User{}
			if tt.expectedUser != emptyUser {
				require.Equal(t, u, tt.expectedUser)
			}
			if tt.expectedError != nil {
				require.Error(t, err, tt.expectedError)
			}

		})
	}
}
