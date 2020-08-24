package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	v1 "sample-order/api/v1"
	"sample-order/api/v1/item"
	"sample-order/config"
	itemCore "sample-order/core/item"
	"sample-order/libs"
	itemRepository "sample-order/repository/item"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

//start server
func start(config *config.AppConfig, db *libs.DatabaseConnection, e *echo.Echo) {
	// run server
	go func() {
		address := fmt.Sprintf("localhost:%d", config.Port)

		if err := e.Start(address); err != nil {
			log.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	//close db
	db.CloseConnection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func main() {
	config := config.InitConfig()

	//load config if available or set to default
	db := libs.NewDatabaseConnection(config)

	var itemDataRepo itemCore.DataRepository
	if db.Driver == libs.MySQL {
		itemDataRepo = itemRepository.NewMySQLRepository(db.MySQLDB)
	} else if db.Driver == libs.MongoDB {
		itemDataRepo = itemRepository.NewMongoDBRepository(db.MongoDB)
	}

	//initiate item service
	itemService := itemCore.NewServiceImpl(itemDataRepo)

	//initiate item controller
	itemControllerV1 := item.NewController(itemService)

	e := echo.New()
	v1.RegisterVIPath(e, itemControllerV1)

	//start server
	start(config, db, e)
}
