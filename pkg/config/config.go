package config

import "os"

type Config struct {
	BotName  string `json:"bot_name"`
	BotToken string `json:"bot_token"`
	ClientId string `json:"client_id"`
}

func Get() *Config {
	return &Config{
		BotName:  os.Getenv("BOT_NAME"),
		BotToken: os.Getenv("BOT_TOKEN"),
		ClientId: os.Getenv("CLIENT_ID"),
	}
}
