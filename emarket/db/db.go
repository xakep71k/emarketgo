package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoClient() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		fmt.Printf("new client: %v\n", err)
	}

	ctx := DefaultContext()
	err = client.Connect(ctx)
	if err != nil {
		fmt.Printf("connect: %v\n", err)
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Printf("ping: %v\n", err)
		return nil, err
	}

	return client, nil
}

func DefaultContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	return ctx
}
