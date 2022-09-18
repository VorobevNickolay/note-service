package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type store interface {
	CreateUser(name, password string) (User, error)
	FindUserByName(name string) (User, error)
}

type Service struct {
	store store
}

func NewService(store store) *Service {
	return &Service{store: store}
}

func (s *Service) SignUp(name, password string) (User, error) {
	password = s.createHash(password)
	user, err := s.store.CreateUser(name, password)
	if err != nil {
		return User{}, fmt.Errorf("failed to signup: %w", err)
	}
	return user, nil
}

func (s *Service) Login(name, password string) (User, error) {
	u, err := s.store.FindUserByName(name)
	if err != nil {
		return User{}, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return User{}, ErrUserNotFound
	}
	return u, nil
}

func (s *Service) createHash(str string) string {
	bytePassword := []byte(str)
	hashPassword, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	return string(hashPassword)
}
