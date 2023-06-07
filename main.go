package main

import (
	"billing/config"
	"billing/db"
	"billing/internal/handler"
	"billing/internal/model"
	"billing/internal/repository"
	"billing/internal/service"
	"billing/logging"
	"github.com/gin-gonic/gin"
	"net"
)

func main() {
	log := logging.GetLogger()

	router := gin.Default()

	cfg := config.GetConfig()

	dbConnection, err := db.GetDBConnection(cfg.Db)
	if err != nil {
		log.Fatal(err)
	}

	if !dbConnection.Migrator().HasTable(&model.Account{}) {
		err := dbConnection.AutoMigrate(&model.Account{})
		if err != nil {
			log.Fatal(err)
		}
	}

	newRepository := repository.NewRepository(dbConnection, log)

	newService := service.NewService(newRepository, log)

	newHandler := handler.NewHandler(router, newService, log)
	newHandler.Init()

	address := net.JoinHostPort(cfg.Listen.Host, cfg.Listen.Port)

	log.Fatal(router.Run(address))
}
