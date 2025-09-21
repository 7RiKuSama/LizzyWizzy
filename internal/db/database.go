package db

import (

	"github.com/7RikuSama/liz.git/internal/events"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	Client *mongo.Client
	Handlers *events.Handlers
}


func NewDatabase(client *mongo.Client) *Database {
	return &Database {
		Client: client,
	}
}
