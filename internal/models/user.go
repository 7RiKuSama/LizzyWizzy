package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID        string             `bson:"user_id" json:"user_id"`
	IsBot         bool               `bson:"-" json:"is_bot"`
}

func NewUser(userID string, IsBot bool) *User {
	return &User {
		UserID: userID,
		IsBot: IsBot,
	}
}
