package service

import (
	"github.com/vavelour/chat/internal/domain/entities"
)

type PrivateRepository interface {
	InsertMessage(m entities.Message) error
	GetMessages(sender, recipient string, limit, offset int) ([]entities.Message, error)
	GetUsers(user string) ([]string, error)
}

//go:generate mockgen -source=private_service.go -destination=mocks/private_repository_mock.go

type PrivateService struct {
	repos PrivateRepository
}

func NewPrivateService(r PrivateRepository) *PrivateService {
	return &PrivateService{repos: r}
}

func (s *PrivateService) SendPrivateMessage(m entities.Message) error {
	return s.repos.InsertMessage(m)
}

func (s *PrivateService) GetPrivateMessages(sender, recipient string, limit, offset int) ([]entities.Message, error) {
	return s.repos.GetMessages(sender, recipient, limit, offset)
}

func (s *PrivateService) ViewUsers(user string) ([]string, error) {
	list, err := s.repos.GetUsers(user)
	if err != nil {
		return nil, err
	}

	return list, nil
}
