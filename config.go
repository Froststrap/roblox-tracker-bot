package main

import (
	"fmt"
	"os"
)

type Config struct {
	BotToken       string
	DiscordGuildId string
}

func loadConfig() (*Config, error) {
	config := &Config{
		BotToken:       os.Getenv("BOT_TOKEN"),
		DiscordGuildId: os.Getenv("DISCORD_GUILD_ID"),
	}

	if config.BotToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN is not set")
	}

	return config, nil
}
