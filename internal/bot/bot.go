package bot

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/7RikuSama/liz.git/internal/db"
	"github.com/7RikuSama/liz.git/internal/events"
	"github.com/7RikuSama/liz.git/internal/service"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session  *discordgo.Session
	Database *db.Database
	Handlers *events.Handlers
}

// Function For creating a bot instance
func NewBot(token string, ctx context.Context) (*Bot, error) {
	if token == "" {
		return nil, errors.New("The token was not provided")
	}

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()

	dbClient := db.DBConnection(ctx)
	services := service.NewServices(dbClient)
	handlers := events.NewHandlers(ctx, services, s)

	return &Bot{
		Session:  s,
		Database: db.NewDatabase(dbClient),
		Handlers: handlers,
	}, nil
}

// Function For opening the bot session
func (b *Bot) Open() error { // ctx is not used, so we remove it

	b.Handlers.RegisterEventHandlers(b.Database.Client, b.Session)
	b.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is ready. Registering commands...")
		b.Handlers.AddCommands(s, "1304848320132415508")
		log.Println("Commands registered successfully.")
	})

	if err := b.Session.Open(); err != nil {
		return err
	}

	// status
	err := b.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: "idle",
		Activities: []*discordgo.Activity{
			{
				Name: "to some fine tunes 🎶",
				Type: discordgo.ActivityTypeListening,
			},
		},
	})

	if err != nil {
		log.Println("Failed to update status:", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	// On shutdown, remove commands
	b.Handlers.DeleteCommands(b.Session, "1304848320132415508")

	log.Println("Removing commands...")

	return b.Stop()
}

// Function for closing the bot session
func (b *Bot) Stop() error {
	return b.Session.Close()
}
