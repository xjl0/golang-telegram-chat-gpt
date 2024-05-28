package tg

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"golang-telegram-chat-gpt/internal/app/gpt"
	"log"
)

type App struct {
	TGBot *bot.Bot
}

func NewApp(token string, gpt *gpt.App) *App {
	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			defaultHandler(ctx, bot, update, gpt)
		}),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/new", bot.MatchTypeExact, func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		newConversationHandler(ctx, bot, update, gpt)
	})

	if _, err := b.SetMyCommands(context.Background(), &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{{
			Command:     "new",
			Description: "New conversation",
		}},
	}); err != nil {
		log.Printf("Failed to set commands: %v", err)
	}

	return &App{TGBot: b}
}
