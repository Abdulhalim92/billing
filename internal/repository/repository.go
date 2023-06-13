package repository

import (
	"billing/internal/model"
	"billing/logging"
	"gorm.io/gorm"
)

type Account interface {
	CreateAccount(account *model.Account) (int, error)
	Transfer(operation *model.Operation) error
}

type Repository struct {
	Account
}

func NewRepository(conn *gorm.DB, logger *logging.Logger) *Repository {
	return &Repository{
		Account: NewAccountRepository(conn, logger),
	}
}
