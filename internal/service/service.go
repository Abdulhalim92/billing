package service

import (
	"billing/internal/repository"
	"billing/logging"
)

type Service struct {
	Repository *repository.Repository
	Logger     *logging.Logger
}

func NewService(rep *repository.Repository, log *logging.Logger) *Service {
	return &Service{
		Repository: rep,
		Logger:     log,
	}
}
