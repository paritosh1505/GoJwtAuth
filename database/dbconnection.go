package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnect() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error while opening the env file")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	mongoUri := os.Getenv("MONGO_URL")
	clientOption := options.Client().ApplyURI(mongoUri)
	connectDb, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		log.Fatal("Error in mongo connection may be env file issue->", err)
	}
	fmt.Println("Mongo db connection in progress...")
	//Check connection

	if err := connectDb.Ping(ctx, nil); err != nil {
		log.Fatal("Ping for mongo db is not happening->", err)
	}
	fmt.Println("Connected to mongo")
	return connectDb
}

var ClientInstance *mongo.Client = MongoConnect()

func CreateCollection(mongoInstance *mongo.Client, TableName string) *mongo.Collection {
	var createInst = mongoInstance.Database("First").Collection(TableName)
	return createInst
}
