package tg

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sashabaranov/go-openai"
	"io"
	"time"
)

type gptProvider interface {
	Send(ctx context.Context, chatID int64, text string) *openai.ChatCompletionStream
	Save(chatID int64, text string)
	Clear(tgID int64)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update, gpt gptProvider) {
	ctxStream, cancel := context.WithTimeout(ctx, time.Minute*2)
	defer cancel()

	stream := gpt.Send(ctxStream, update.Message.Chat.ID, update.Message.Text)
	if stream == nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Sorry, something went wrong. Please try again later.",
		}); err != nil {
			fmt.Printf("\nSend message error: %v\n", err)
		}
		return
	}
	defer func(stream *openai.ChatCompletionStream) {
		if err := stream.Close(); err != nil {
			fmt.Printf("\nStream close error: %v\n", err)
		}
	}(stream)

	throttleTime := time.Now()
	content := ""
	msgID := 0

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Sorry, something went wrong. Please try again later.",
			}); err != nil {
				fmt.Printf("\nSend message error: %v\n", err)
			}
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		if response.Choices[0].Delta.Content == "" {
			continue
		}

		content += response.Choices[0].Delta.Content

		if msgID == 0 {
			msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   content,
			})
			if err != nil {
				fmt.Printf("\nSend message error: %v\n", err)
				return
			}
			msgID = msg.ID
		} else {
			if throttleTime.Add(time.Millisecond * 1300).Before(time.Now()) {
				if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
					ChatID:    update.Message.Chat.ID,
					Text:      content,
					MessageID: msgID,
				}); err != nil {
					fmt.Printf("\nEdit message error: %v\n", err)
				}

				throttleTime = time.Now()
			}
		}
	}

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.Message.Chat.ID,
		Text:      content,
		MessageID: msgID,
	}); err != nil {
		fmt.Printf("\nEdit message error: %v\n", err)
	}

	gpt.Save(update.Message.Chat.ID, content)
}
func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hi, I'm ChatGPT. Ask me anything!",
	}); err != nil {
		fmt.Printf("\nSend message error: %v\n", err)
	}
}

func newConversationHandler(ctx context.Context, b *bot.Bot, update *models.Update, gpt gptProvider) {
	gpt.Clear(update.Message.Chat.ID)
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Conversation cleared.",
	}); err != nil {
		fmt.Printf("\nSend message error: %v\n", err)
	}
}
