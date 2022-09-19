package user

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetUsers(t *testing.T) {
	t.Run("should return empty list", func(t *testing.T) {
		store := NewInMemoryStore()
		actual, err := store.GetUsers()
		require.NoError(t, err)
		require.Equal(t, 0, len(actual))
	})

	t.Run("should return users", func(t *testing.T) {
		store := NewInMemoryStore()
		exp1, err := store.CreateUser(uuid.NewString(), uuid.NewString())
		require.NoError(t, err)

		exp2, err := store.CreateUser(uuid.NewString(), uuid.NewString())
		require.NoError(t, err)

		actual, err := store.GetUsers()
		require.NoError(t, err)

		require.Equal(t, 2, len(actual))
		if exp1 == *actual[0] {
			require.Equal(t, exp1, *actual[0])
			require.Equal(t, exp2, *actual[1])
		} else {
			require.Equal(t, exp1, *actual[1])
			require.Equal(t, exp2, *actual[0])
		}
	})
}

func TestFindUserByID(t *testing.T) {
	t.Run("should return ErrUserNotFound", func(t *testing.T) {
		store := NewInMemoryStore()

		actual, err := store.FindUserByID(uuid.NewString())
		require.Error(t, err, ErrUserNotFound)
		require.Equal(t, actual, User{})
	})

	t.Run("should find user", func(t *testing.T) {
		store := NewInMemoryStore()

		_, err := store.CreateUser(uuid.NewString(), uuid.NewString())
		require.NoError(t, err)

		expected, err := store.CreateUser(uuid.NewString(), uuid.NewString())
		require.NoError(t, err)

		_, err = store.CreateUser(uuid.NewString(), uuid.NewString())
		require.NoError(t, err)

		actual, err := store.FindUserByID(expected.ID)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}

func TestFindUserByName(t *testing.T) {
	t.Run("should return ErrUserNotFound", func(t *testing.T) {
		store := NewInMemoryStore()

		actual, err := store.findUserByName(uuid.NewString())
		require.Error(t, err, ErrUserNotFound)
		require.Equal(t, actual, User{})
	})

	t.Run("should find user", func(t *testing.T) {
		store := NewInMemoryStore()

		_, err := store.CreateUser(uuid.NewString(), uuid.NewString())
		require.NoError(t, err)

		expected, err := store.CreateUser(uuid.NewString(), uuid.NewString())
		require.NoError(t, err)

		_, err = store.CreateUser(uuid.NewString(), uuid.NewString())
		require.NoError(t, err)

		actual, err := store.findUserByName(expected.Username)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
	t.Run("should return ErrUserNotFound", func(t *testing.T) {
		store := NewInMemoryStore()

		actual, err := store.FindUserByName(uuid.NewString())
		require.Error(t, err, ErrUserNotFound)
		require.Equal(t, actual, User{})
	})

	t.Run("should find user", func(t *testing.T) {
		store := NewInMemoryStore()

		_, err := store.CreateUser(uuid.NewString(), uuid.NewString())
		require.NoError(t, err)

		expected, err := store.CreateUser(uuid.NewString(), uuid.NewString())
		require.NoError(t, err)

		_, err = store.CreateUser(uuid.NewString(), uuid.NewString())
		require.NoError(t, err)

		actual, err := store.FindUserByName(expected.Username)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("should return errUsedUserName", func(t *testing.T) {
		store := NewInMemoryStore()

		username := uuid.NewString()
		_, err1 := store.CreateUser(username, uuid.NewString())
		actual, err2 := store.CreateUser(username, uuid.NewString())
		require.NoError(t, err1)
		require.Error(t, err2, ErrUsedUsername)
		require.Equal(t, actual, User{})
	})
	t.Run("should create user", func(t *testing.T) {
		store := NewInMemoryStore()

		expected := User{
			Username: uuid.NewString(),
			ID:       uuid.NewString(),
			Password: uuid.NewString(),
		}
		actual, err := store.CreateUser(expected.Username, expected.Password)
		require.NoError(t, err)
		require.Equal(t, expected.Username, actual.Username)
		require.Equal(t, expected.Password, actual.Password)
		require.NotEmpty(t, actual.ID)
	})
}
