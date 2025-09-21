package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Membership struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID   string             `bson:"user_id,omitempty" json:"user_id"`
	GuildID  string             `bson:"guild_id,omitempty" json:"guild_id"`
	Nick     string             `bson:"nickname" json:"nickname"`
	Eddies   int                `bson:"eddies" json:"eddies"`
	Exp      int                `bson:"exp" json:"exp"`
	Level    int                `bson:"level" json:"level"`
	JoinedAt time.Time          `bson:"joined_at" json:"joined_at"`
}

func NewMembership(userID, guildID, Nick string) *Membership {
	return &Membership{
		UserID:   userID,
		GuildID:  guildID,
		Nick:     Nick,
		Exp:      0,
		Level:    1,
		Eddies:   250,
		JoinedAt: time.Now(),
	}
}
