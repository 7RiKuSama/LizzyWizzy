package events

import (
	"context"

	"github.com/7RikuSama/liz.git/internal/commands"
	"github.com/7RikuSama/liz.git/internal/service"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Player *commands.Player
)

type Handlers struct {
	Session  *discordgo.Session
	Services *service.Services
	Context  context.Context
	Slash    *Slash
}

func NewHandlers(ctx context.Context, services *service.Services, s *discordgo.Session) *Handlers {
	h := &Handlers{
		Session:  s,
		Services: services,
		Context:  ctx,
		Slash:    NewSlash(),
	}
	h.initCommands()
	return h
}

func (h *Handlers) RegisterEventHandlers(client *mongo.Client, s *discordgo.Session) {
	if Player == nil {
		Player = commands.NewPlayer()
		Player.AddMusicFiles()
	}
	s.AddHandler(h.OnInteractionCreate)
	s.AddHandler(h.OnMessageCreate)
}
