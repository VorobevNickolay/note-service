package user

import (
	"context"
	"github.com/google/uuid"
	"strings"
	"sync"
)

type inMemoryStore struct {
	sync.RWMutex
	users map[string]User
}

func NewInMemoryStore() *inMemoryStore {
	return &inMemoryStore{users: make(map[string]User)}
}

func (store *inMemoryStore) CreateUser(ctx context.Context, name, password string) (User, error) {
	store.Lock()
	defer store.Unlock()

	if len(password) == 0 || len(name) == 0 {
		return User{}, ErrEmptyPassword
	}
	if _, err := store.findUserByName(ctx, name); err == nil {
		return User{}, ErrUsedUsername
	}
	user := User{
		ID:       uuid.NewString(),
		Username: name,
		Password: password,
	}
	store.users[user.ID] = user
	return user, nil
}

func (store *inMemoryStore) FindUserById(_ context.Context, id string) (User, error) {
	store.RLock()
	defer store.RUnlock()

	if u, ok := store.users[id]; ok {
		return u, nil
	}
	return User{}, ErrUserNotFound
}

func createPointer(u User) *User {
	return &u
}

func (store *inMemoryStore) GetUsers(_ context.Context) ([]*User, error) {
	store.RLock()
	defer store.RUnlock()
	res := make([]*User, len(store.users))
	i := 0

	for j := range store.users {
		res[i] = createPointer(store.users[j])
		i++
	}

	return res, nil
}

// findUserByName find user and isn't thread-safe
func (store *inMemoryStore) findUserByName(_ context.Context, name string) (User, error) {
	for _, u := range store.users {
		if strings.EqualFold(name, u.Username) {
			return u, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (store *inMemoryStore) FindUserByNameAndPassword(ctx context.Context, name, password string) (User, error) {
	store.RLock()
	defer store.RUnlock()

	u, err := store.findUserByName(ctx, name)
	if err != nil {
		return User{}, err
	}
	if password != u.Password {
		return User{}, ErrUserNotFound
	}
	return u, nil
}
