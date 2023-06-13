package service

import (
	"billing/internal/model"
	"billing/internal/repository"
	"billing/logging"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Account interface {
	CreateAccount(account *model.Account) (int, error)
	Transfer(operation *model.Operation) error
}

type Service struct {
	Account
}

func NewService(rep *repository.Repository, log *logging.Logger) *Service {
	return &Service{
		Account: NewAccountService(rep, log),
	}
}
