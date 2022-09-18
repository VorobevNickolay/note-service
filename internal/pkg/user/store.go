package user

import (
	"strings"
	"sync"

	"github.com/google/uuid"
)

type inMemoryStore struct {
	sync.RWMutex
	users map[string]User
}

func NewInMemoryStore() *inMemoryStore {
	return &inMemoryStore{users: make(map[string]User)}
}

func (store *inMemoryStore) CreateUser(name, password string) (User, error) {
	store.Lock()
	defer store.Unlock()

	if _, err := store.findUserByName(name); err == nil {
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

func (store *inMemoryStore) FindUserByID(id string) (User, error) {
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

func (store *inMemoryStore) GetUsers() ([]*User, error) {
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

func (store *inMemoryStore) FindUserByName(name string) (User, error) {
	store.RLock()
	defer store.RUnlock()
	return store.findUserByName(name)
}

// findUserByName find user and isn't thread-safe
func (store *inMemoryStore) findUserByName(name string) (User, error) {
	for _, u := range store.users {
		if strings.EqualFold(name, u.Username) {
			return u, nil
		}
	}
	return User{}, ErrUserNotFound
}
