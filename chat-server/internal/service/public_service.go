package service

import "github.com/vavelour/chat/internal/domain/entities"

//go:generate mockgen -source=public_service.go -destination=mocks/public_repository_mock.go

type PublicRepository interface {
	InsertMessage(m entities.Message) error
	GetMessages(limit, offset int) ([]entities.Message, error)
}

type PublicService struct {
	repos PublicRepository
}

func NewPublicService(r PublicRepository) *PublicService {
	return &PublicService{repos: r}
}

func (s *PublicService) SendPublicMessage(m entities.Message) error {
	return s.repos.InsertMessage(m)
}

func (s *PublicService) GetPublicMessages(limit, offset int) ([]entities.Message, error) {
	return s.repos.GetMessages(limit, offset)
}
