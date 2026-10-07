package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Info("no .env file found, using environment variables")
	}

	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	slog.SetDefault(logger)

	config, err := loadConfig()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	bot, err := newBot(config.BotToken)
	if err != nil {
		slog.Error(
			"failed to create Discord bot",
			"error",
			err,
		)
		os.Exit(1)
	}

	if err := bot.start(); err != nil {
		slog.Error(
			"failed to start Discord bot",
			"error",
			err,
		)
		os.Exit(1)
	}

	slog.Info("configuration loaded")
	slog.Info("bot is running")

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()

	slog.Info("shutting down")

	if err := bot.close(); err != nil {
		slog.Error("failed to close Discord bot", "error", err)
	}
}
