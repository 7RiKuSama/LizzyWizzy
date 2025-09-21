package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Guild struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GuildID string             `bson:"guild_id,omitempty" json:"guild_id"`
	Prefix       string             `bson:"prefix,omitempty" json:"prefix"`
	BannedWords  []string           `bson:"banned_words,omitempty" json:"banned_words"`
	AdminRole    string             `bson:"admin_role,omitempty" json:"admin_role"`
	CustomStatus string             `bson:"custom_status,omitempty" json:"custom_status"`
	CustomName   string             `bson:"custom_name,omitempty" json:"custom_name"`
	WelcomeChan  string             `bson:"welcome_chan,omitempty" json:"welcome_chan"`
	LogChan      string             `bson:"log_chan,omitempty" json:"log_chan"`
	MusicChan    string             `bson:"music_chan,omitempty" json:"music_chan"`
}


func NewGuild(guildID string) *Guild {
	return &Guild {
		GuildID: guildID,
	}
}
