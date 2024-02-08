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
}

func NewService(config *config.Config, repository *repository.Repository) ServiceImpl {
	s := &Service{
		config:     config,
		repository: repository,
	}

	return s
}
