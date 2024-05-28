package app

import (
	"golang-telegram-chat-gpt/internal/app/gpt"
	"golang-telegram-chat-gpt/internal/app/tg"
	"golang-telegram-chat-gpt/internal/config"
)

type App struct {
	TGBot *tg.App
}

func NewApp(cfg config.Config) *App {

	chat := gpt.NewApp(cfg.OpenApiToken, cfg.Proxy, cfg.ConversationIdleTimeOut, cfg.Model)

	return &App{TGBot: tg.NewApp(cfg.TelegramToken, chat)}
}
