package handler

import (
	"billing/internal/model"
	"billing/internal/service"
	"billing/logging"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Handler struct {
	Engine  *gin.Engine
	Service *service.Service
	Logger  *logging.Logger
}

func NewHandler(engine *gin.Engine, service *service.Service, logger *logging.Logger) *Handler {
	return &Handler{
		Engine:  engine,
		Service: service,
		Logger:  logger,
	}
}

func (h *Handler) Init() {
	router := h.Engine.Group("v1")

	router.POST("/transfer", h.Transfer)
	router.POST("/create", h.CreateAccount)
}

func (h *Handler) CreateAccount(c *gin.Context) {
	var account *model.Account

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

	if err = h.Service.CreateAccount(account); err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, "error creating account")
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, "Successfully create account")
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
