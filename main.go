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

	"go.mongodb.org/mongo-driver/mongo/readpref"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type appConfig struct {
	Port     int `yaml:"port"`
	Database struct {
		Driver   string `yaml:"driver"`
		Name     string `yaml:"name"`
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
}

func initConfig() *appConfig {
	var defaultConfig appConfig
	defaultConfig.Port = 1323
	defaultConfig.Database.Driver = "mongodb"
	defaultConfig.Database.Name = "transaction"
	defaultConfig.Database.Address = "localhost"
	defaultConfig.Database.Port = 27017
	defaultConfig.Database.Username = ""
	defaultConfig.Database.Password = ""

	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("config file not found, will use default value")
		} else {
			log.Info("error to load config file, will use default value")
		}

		return &defaultConfig
	}

	var finalConfig appConfig
	err := viper.Unmarshal(&finalConfig)
	if err != nil {
		log.Info("failed to extract config, will use default value")
		return &defaultConfig
	}

	return &finalConfig
}

func newMysqlDB(config *appConfig) *sql.DB {
	var uri string

	uri = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		config.Database.Username,
		config.Database.Password,
		config.Database.Address,
		config.Database.Port,
		config.Database.Name)

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

func newMongoDBClient(config *appConfig) *mongo.Client {
	uri := "mongodb://"

	if config.Database.Username != "" {
		uri = fmt.Sprintf("%s%v:%v@", uri, config.Database.Username, config.Database.Password)
	}

	uri = fmt.Sprintf("%s%v:%v",
		uri,
		config.Database.Address,
		config.Database.Port)

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
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
	config := initConfig()

	if config.Database.Driver == "mysql" {
		//initiate mysql db repository
		db := newMysqlDB(config)
		defer db.Close()
		itemRepo = itemDomain.NewMySQLRepository(db)
	} else if config.Database.Driver == "mongodb" {
		// //initiate mongodb repository
		client := newMongoDBClient(config)
		defer client.Disconnect(context.Background())
		db := client.Database(config.Database.Name)
		itemRepo = itemDomain.NewMongoDBRepository(db)
	} else {
		panic("Unsupported database driver")
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
