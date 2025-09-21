package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/7RikuSama/liz.git/internal/bot"
	"github.com/joho/godotenv"
)

var botInstance *bot.Bot

func main() {
	// hello world 
	godotenv.Load()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	
	log.Println("Bot started. if you want to stop it press 'CTRL+C'")
	botInstance, err := bot.LaunchBot(os.Getenv("BOT_TOKEN_2"), ctx) 

	if err != nil {
		panic(err)
	}

	defer botInstance.Stop()
}
