package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	businessItem "sample-order/business/item"
	"sample-order/config"
	api "sample-order/modules/api"
	itemControllerV1 "sample-order/modules/api/v1/item"
	itemRepo "sample-order/modules/repository/item"
	"sample-order/util"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func main() {
	config := config.GetConfig()

	//load config if available or set to default
	dbCon := util.NewDatabaseConnection(config)

	//initiate item repository
	itemRepo := itemRepo.RepositoryFactory(dbCon)

	//initiate item service
	itemService := businessItem.NewService(itemRepo)

	//initiate item controller
	itemControllerV1 := itemControllerV1.NewController(itemService)

	//create echo http
	e := echo.New()

	//register API path and handler
	api.RegisterPath(e, itemControllerV1)

	// run server
	go func() {
		address := fmt.Sprintf("localhost:%d", config.Port)

		if err := e.Start(address); err != nil {
			log.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	//close db
	defer dbCon.CloseConnection()

	// a timeout of 10 seconds to shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
