package service

import (
	"errors"
	"github.com/vavelour/chat/internal/domain/entities"
)

var (
	ErrIncorrectPassword       = errors.New("incorrect password")
	ErrIncorrectTypeConversion = errors.New("incorrect type conversion")
)

//go:generate mockgen -source=basic_auth_service.go -destination=mocks/auth_repository_mock.go

type AuthRepository interface {
	InsertUser(username, password string) error
	GetUser(username string) (entities.User, error)
}

type AuthService struct {
	repos AuthRepository
}

func NewAuthService(r AuthRepository) *AuthService {
	return &AuthService{repos: r}
}

func (s *AuthService) CreateUser(username, password string) (string, error) {
	if err := s.repos.InsertUser(username, password); err != nil {
		return "", err
	}

	return username, nil
}

func (s *AuthService) UserIdentity(usr interface{}) (string, error) {
	u, ok := usr.(entities.User)
	if !ok {
		return "", ErrIncorrectTypeConversion
	}

	user, err := s.repos.GetUser(u.Username)
	if err != nil {
		return "", err
	}

	if user.Password != u.Password {
		return "", ErrIncorrectPassword
	}

	return u.Username, nil
}
