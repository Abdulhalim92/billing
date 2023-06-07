package repository

import (
	"billing/logging"
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
