package service

import (
	"billing/internal/model"
	"billing/internal/repository"
	"billing/logging"
)

type AccountService struct {
	repository *repository.Repository
	logger     *logging.Logger
}

func NewAccountService(repository *repository.Repository, logger *logging.Logger) *AccountService {
	return &AccountService{
		repository: repository,
		logger:     logger,
	}
}

func (s *AccountService) CreateAccount(account *model.Account) (int, error) {
	return s.repository.CreateAccount(account)
}

func (s *AccountService) Transfer(operation *model.Operation) error {
	return s.repository.Transfer(operation)
}
