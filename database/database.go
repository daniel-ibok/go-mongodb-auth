package database

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoDBURI        string
	MongoDBName       string
	MongoDBCollection string
}

var config *Config
var DB *mongo.Database
var Client *mongo.Client

func init() {

	config = &Config{
		MongoDBURI:        os.Getenv("MONGODB_URI"),
		MongoDBName:       os.Getenv("MONGODB_DB_NAME"),
		MongoDBCollection: os.Getenv("MONGODB_COLLECTION"),
	}

}

func GetDBCollection() *mongo.Collection {
	return DB.Collection(config.MongoDBCollection)
}

func NewDBInstance() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.MongoDBURI).SetServerAPIOptions(serverAPI))
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database(config.MongoDBName)
	Client = client
	return nil
}

func Close() error {
	return Client.Disconnect(context.Background())
}
