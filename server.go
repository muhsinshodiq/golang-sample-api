package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	itemAPIV1 "sample-order/api/v1/item"
	itemDomain "sample-order/domain/item"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newConfig() {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")

	database := make(map[string]interface{})
	database["driver"] = "mongodb"
	database["address"] = "localhost"
	database["port"] = 27017
	database["username"] = ""
	database["password"] = ""
	database["name"] = "transaction"
	viper.SetDefault("database", database)

	server := make(map[string]interface{})
	server["port"] = 1323
	viper.SetDefault("port", 1323)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("config file not found, will use default value")
		} else {
			log.Info("error to load config file, will use default value")
		}
	}
}

func newMysqlDB() *sql.DB {
	var uri string

	database := viper.GetStringMap("database")

	uri = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		database["username"],
		database["password"],
		database["address"],
		database["port"],
		database["name"])

	db, err := sql.Open("mysql", uri)
	if err != nil {
		log.Info("failed to connect database: ", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Info("failed to connect database: ", err)
		panic(err)
	}

	return db
}

func newMongoDBClient() *mongo.Client {
	uri := "mongodb://"

	database := viper.GetStringMap("database")

	if database["username"] != "" {
		uri = fmt.Sprintf("%s%v:%v@", uri, database["username"], database["password"])
	}

	uri = fmt.Sprintf("%s%v:%v",
		uri,
		database["address"],
		database["port"])

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	return client
}

func startServer(e *echo.Echo) {
	// Start server
	go func() {
		address := fmt.Sprintf("localhost:%d", viper.GetInt("port"))

		if err := e.Start(address); err != nil {
			log.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var itemRepo itemDomain.Repository

	//load config if available or set to default
	newConfig()

	database := viper.GetStringMap("database")

	if database["driver"].(string) == "mysql" {
		//initiate mysql db repository
		db := newMysqlDB()
		defer db.Close()
		itemRepo = itemDomain.NewMySQLRepository(db)
	} else {
		// //initiate mongodb repository
		client := newMongoDBClient()
		defer client.Disconnect(context.Background())
		db := client.Database(database["name"].(string))
		itemRepo = itemDomain.NewMongoDBRepository(db)
	}

	//initiate item service
	itemService := itemDomain.NewServiceImpl(itemRepo)

	//initiate item controller
	itemControllerV1 := itemAPIV1.NewController(itemService)

	e := echo.New()

	//register item controller v1 path
	itemControllerV1.RegisterPath(e)

	startServer(e)
}
