package main

import (
	"context"
	"golang-telegram-chat-gpt/internal/app"
	"golang-telegram-chat-gpt/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg := config.MustLoad()

	application := app.NewApp(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	application.TGBot.TGBot.Start(ctx)

}
