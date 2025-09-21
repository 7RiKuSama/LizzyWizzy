package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBConnection(ctx context.Context) *mongo.Client {
	uri := "mongodb+srv://rikusama:UfrXULGFageF8KcC@cluster0.havtds4.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Panic(err)
	}
	return client
}

func DBDisconnect(ctx context.Context, client *mongo.Client) {
	if err := client.Disconnect(ctx); err != nil {
		log.Panic(err)
	}
}
