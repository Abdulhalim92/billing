package handler

import (
	"billing/internal/model"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (h *Handler) CreateAccount(c *gin.Context) {
	var account model.Account

	err := c.BindJSON(&account)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, "please send valid data")
		c.Abort()
		return
	}

	if account.Balance < 0 {
		h.Logger.Error("please send valid data")
		c.JSON(http.StatusBadRequest, "please send valid data")
		c.Abort()
		return
	}

	if account.Status == "" {
		h.Logger.Error("please send valid data")
		c.JSON(http.StatusBadRequest, "please send valid data")
		c.Abort()
		return
	}

	id, err := h.Service.CreateAccount(&account)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, "something went wrong")
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully create account", "ID": id})
}

func (h *Handler) Transfer(c *gin.Context) {
	time.Sleep(1 * time.Second)

	var operation *model.Operation

	err := c.BindJSON(&operation)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, "please send valid data")
		c.Abort()
		return
	}

	if operation.AccountFromID <= 0 || operation.AccountToID <= 0 && operation.Amount <= 0 {
		h.Logger.Error(errors.New("access denied for negative numbers in request body"))
		c.JSON(http.StatusBadRequest, "access denied for negative numbers in request body")
		c.Abort()
		return
	}

	if err = h.Service.Transfer(operation); err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, "error occurred while committing transaction")
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, "Successfully transaction")
}
