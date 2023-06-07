package service

import (
	"billing/internal/model"
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

func (s *Service) CreateAccount(account *model.Account) error {
	return s.Repository.CreateAccount(account)
}

func (s *Service) Transfer(operation *model.Operation) error {
	return s.Repository.Transfer(operation)
}
