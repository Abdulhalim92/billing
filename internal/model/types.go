package model

import (
	"time"
)

type Account struct {
	ID      int       `json:"-"`
	Number  string    `json:"name"`
	Balance float64   `json:"balance" binding:"required" gorm:"type:numeric(10, 2);size:10"`
	Status  string    `json:"status"`
	Created time.Time `json:"-"`
	Updated time.Time `json:"-"`
}

type Operation struct {
	AccountFromID int     `json:"account_from" binding:"required"`
	AccountToID   int     `json:"account_to" binding:"required"`
	Amount        float64 `json:"amount" binding:"required" gorm:"type:numeric(10, 2);size:10"`
}
