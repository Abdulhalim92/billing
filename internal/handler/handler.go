package handler

import (
	"billing/internal/service"
	"billing/logging"
	"github.com/gin-gonic/gin"
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
