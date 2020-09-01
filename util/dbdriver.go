package util

import (
	"context"
	"database/sql"
	"fmt"
	"sample-order/config"
	"time"

	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//DatabaseDriver Database driver enum
type DatabaseDriver string

const (
	//MongoDB MongoDB DatabaseDriver
	MongoDB DatabaseDriver = "mongodb"
	//MySQL MySQL DatabaseDriver
	MySQL DatabaseDriver = "mysql"
)

//DatabaseConnection Database connection
type DatabaseConnection struct {
	Driver DatabaseDriver

	//for MySQL
	MySQLDB *sql.DB

	//for MongoDB
	MongoDB     *mongo.Database
	mongoClient *mongo.Client
}

//NewDatabaseConnection Create new database connection based on given config
func NewDatabaseConnection(config *config.AppConfig) *DatabaseConnection {
	var db DatabaseConnection
	//define the data repository
	if config.Database.Driver == "mysql" {
		//initiate mysql db repository
		db.MySQLDB = newMysqlDB(config)
		db.Driver = MySQL
	} else if config.Database.Driver == "mongodb" {
		// //initiate mongodb repository
		db.mongoClient = newMongoDBClient(config)
		db.MongoDB = db.mongoClient.Database(config.Database.Name)
		db.Driver = MongoDB
	} else {
		panic("Unsupported database driver")
	}

	return &db
}

//CloseConnection Close db connection
func (db *DatabaseConnection) CloseConnection() {
	if db.MySQLDB != nil {
		db.MySQLDB.Close()
	}

	if db.mongoClient != nil {
		db.mongoClient.Disconnect(context.Background())
	}
}

func newMysqlDB(config *config.AppConfig) *sql.DB {
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

func newMongoDBClient(config *config.AppConfig) *mongo.Client {
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
