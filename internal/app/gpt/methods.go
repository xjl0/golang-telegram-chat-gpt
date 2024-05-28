package gpt

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"log"
	"time"
)

// Send message to ChatGPT
func (a *App) Send(ctx context.Context, tgID int64, text string) *openai.ChatCompletionStream {
	history, exists := a.history[tgID]
	if exists && history.t.Add(time.Duration(a.cTimeout)*time.Second).Before(time.Now()) {
		delete(a.history, tgID)
	}

	if history == nil {
		history = &state{
			t:       time.Now(),
			message: make([]openai.ChatCompletionMessage, 0),
		}
		a.history[tgID] = history
	}

	history.message = append(history.message, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: text,
	})

	stream, err := a.client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:    a.model,
			Messages: history.message,
			Stream:   true,
			TopP:     1,
			N:        1,
		},
	)
	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		return nil
	}

	return stream
}

// Save history message
func (a *App) Save(tgID int64, text string) {
	history, exists := a.history[tgID]
	if !exists {
		return
	}
	history.message = append(history.message, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: text,
	})
}

// Clear history
func (a *App) Clear(tgID int64) {
	if _, exists := a.history[tgID]; exists {
		delete(a.history, tgID)
	}
}
