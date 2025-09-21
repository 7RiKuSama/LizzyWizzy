package bot

import (
	"context"
	"os"
)

func LaunchBot(token string, ctx context.Context) (*Bot, error) {
	if token == "" {
		token = os.Getenv("BOT_TOKEN")
	}
	
	bot, err := NewBot(token, ctx)

	if err != nil {
		return nil, err
	}

	if err := bot.Open(); err != nil {
		return nil, err
	}
	return bot, nil
}
