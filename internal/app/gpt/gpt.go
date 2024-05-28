package gpt

import (
	"github.com/sashabaranov/go-openai"
	"net/http"
	"net/url"
	"time"
)

type state struct {
	t       time.Time
	message []openai.ChatCompletionMessage
}

type App struct {
	client   *openai.Client
	history  map[int64]*state
	cTimeout int
	model    string
}

func NewApp(token, proxy string, cTimeout int, model string) *App {
	history := make(map[int64]*state)

	if proxy == "" {
		return &App{client: openai.NewClient(token), history: history, cTimeout: cTimeout, model: model}
	}
	config := openai.DefaultConfig(token)
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}

	return &App{client: openai.NewClientWithConfig(config), history: history, cTimeout: cTimeout, model: model}
}
