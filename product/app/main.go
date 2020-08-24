package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	v1 "sample-order/api/v1"
	"sample-order/config"
	"sample-order/libs"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func main() {
	config := config.GetConfig()

	//load config if available or set to default
	dbCon := libs.NewDatabaseConnection(config)

	//initiate item controller
	itemControllerV1 := InitializeItemControllerV1(dbCon)

	//create echo http
	e := echo.New()

	//register v1 API handler
	v1.RegisterVIPath(e, itemControllerV1)

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
