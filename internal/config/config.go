package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	OpenApiToken            string  `env:"TI_OPENAI_API_KEY" env-required:"true"`
	TelegramToken           string  `env:"TI_TELEGRAM_BOT_TOKEN" env-required:"true"`
	ConversationIdleTimeOut int     `env:"TI_CONVERSATION_IDLE_TIMEOUT" env-default:"3600"`
	Proxy                   string  `env:"TI_PROXY" env-default:""`
	AllowedTGIDs            []int64 `env:"TI_ALLOWED_TG_IDS" env-required:"true"`
	Model                   string  `env:"TI_MODEL" env-default:"gpt-3.5-turbo"`
}

func MustLoad() Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
