package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	Client = client
	return nil
}

func PolicyCollection(dbName string) *mongo.Collection {
	return Client.Database(dbName).Collection("policies")
}

func TopicCollection(dbName string) *mongo.Collection {
	return Client.Database(dbName).Collection("topics")
}
