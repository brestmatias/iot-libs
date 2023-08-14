package wrappers

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClientWrapper struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func (MongoClientWrapper) New(ctx context.Context, uri string, database string) *MongoClientWrapper {
	client := buildMongoClient(uri, ctx)
	return &MongoClientWrapper{
		Client:   client,
		Database: client.Database(database),
	}
}

func (m MongoClientWrapper) GetDatabase() *mongo.Database {
	return m.Database
}

func buildMongoClient(uri string, ctx context.Context) *mongo.Client {
	method := "buildMongoClient"
	log.Printf("[method:%v]🏗️ 🏗️ Building", method)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
