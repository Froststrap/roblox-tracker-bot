package main

import (
	"fmt"
	"os"
)

type Config struct {
	BotToken string
}

func loadConfig() (*Config, error) {
	config := &Config{
		BotToken: os.Getenv("BOT_TOKEN"),
	}

	if config.BotToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN is not set")
	}

	return config, nil
}
