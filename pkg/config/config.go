package config

import "os"

type Config struct {
	BotName   string `json:"bot_name"`
	BotToken  string `json:"bot_token"`
	ClientId  string `json:"client_id"`
	ChannelId string `json:"channel_id"`
	SecretKey string `json:"secret_key"`
}

var GlobalConfig *Config

func Get() *Config {
	if GlobalConfig != nil {
		return GlobalConfig
	} else {
		GlobalConfig := &Config{
			BotName:   os.Getenv("BOT_NAME"),
			BotToken:  os.Getenv("BOT_TOKEN"),
			ClientId:  os.Getenv("CLIENT_ID"),
			ChannelId: os.Getenv("CHANNEL_ID"),
			SecretKey: os.Getenv("SECRET_KEY"),
		}
		return GlobalConfig
	}
}
