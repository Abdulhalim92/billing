package repository

import (
	"billing/internal/model"
	"billing/logging"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Repository struct {
	Connection *gorm.DB
	Logger     *logging.Logger
}

func NewRepository(conn *gorm.DB, logger *logging.Logger) *Repository {
	return &Repository{
		Connection: conn,
		Logger:     logger,
	}
}

func (r *Repository) CreateAccount(account *model.Account) error {
	if err := r.Connection.Create(&account).Error; err != nil {
		r.Logger.Error(err)
		return err
	}

	return nil
}

func (r *Repository) Transfer(operation *model.Operation) error {
	var (
		accountFrom *model.Account
		accountTo   *model.Account
	)

	tx := r.Connection.Begin()

	// получение аккаунтов по ID
	if err := tx.Where("id = ?", operation.AccountFromID).Find(&accountFrom).Error; err != nil {
		r.Logger.Error(err)
		tx.Rollback()
		return err
	}
	if accountFrom.ID == 0 {
		r.Logger.Errorf("no account with %v ID", operation.AccountFromID)
		tx.Rollback()
		return errors.New(fmt.Sprintf("no account with %v ID", operation.AccountFromID))
	}

	if err := tx.Where("id = ?", operation.AccountToID).Find(&accountTo).Error; err != nil {
		r.Logger.Error(err)
		tx.Rollback()
		return err
	}
	if accountFrom.ID == 0 {
		r.Logger.Errorf("no account with %v ID", operation.AccountToID)
		tx.Rollback()
		return errors.New(fmt.Sprintf("no account with %v ID", operation.AccountToID))
	}

	// совершение операций
	// дебит при активном балансе
	if accountFrom.Status == "active" {
		accountFrom.Balance += operation.Amount
		if err := tx.Save(&accountFrom).Error; err != nil {
			r.Logger.Error(err)
			tx.Rollback()
			return err
		}
	}
	// дебит при пассивном балансе
	if accountFrom.Status == "passive" {
		if accountFrom.Balance == 0 || accountFrom.Balance < operation.Amount {
			r.Logger.Error("insufficient balance")
			tx.Rollback()
			return errors.New("insufficient balance")
		}
		accountFrom.Balance -= operation.Amount
		if err := tx.Save(&accountFrom).Error; err != nil {
			r.Logger.Error(err)
			tx.Rollback()
			return err
		}
	}
	// кредит при активном балансе
	if accountTo.Status == "active" {
		if accountTo.Balance == 0 || accountTo.Balance < operation.Amount {
			r.Logger.Error("insufficient balance")
			tx.Rollback()
			return errors.New("insufficient balance")
		}
		accountTo.Balance -= operation.Amount
		if err := tx.Save(&accountTo).Error; err != nil {
			r.Logger.Error(err)
			tx.Rollback()
			return err
		}
	}
	// кредит при пассивном балансе
	if accountTo.Status == "passive" {
		accountTo.Balance += operation.Amount
		if err := tx.Save(&accountTo).Error; err != nil {
			r.Logger.Error(err)
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		r.Logger.Error(err)
		tx.Rollback()
		return err
	}

	return nil
}
