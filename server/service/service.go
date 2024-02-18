package service

import (
	"github.com/04Akaps/kafka-go/config"
	"github.com/04Akaps/kafka-go/server/repository"
)

type Service struct {
	config *config.Config

	repository *repository.Repository
}

type ServiceImpl interface {
	Like(fromUser, toUser string) error
	UnLike(fromUser, toUser string) error
}

func NewService(config *config.Config, repository *repository.Repository) ServiceImpl {
	s := &Service{
		config:     config,
		repository: repository,
	}

	return s
}

func (s *Service) Like(fromUser, toUser string) error {

	return nil
}

func (s *Service) UnLike(fromUser, toUser string) error {
	return nil
}
