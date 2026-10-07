package main

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
}

func newBot(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf(
			"create Discord session: %w",
			err,
		)
	}

	bot := &Bot{
		session: session,
	}

	session.AddHandler(bot.handleReady)

	return bot, nil
}

func (bot *Bot) start() error {
	if err := bot.session.Open(); err != nil {
		return fmt.Errorf(
			"open Discord gateway: %w",
			err,
		)
	}

	return nil
}

func (bot *Bot) close() error {
	if err := bot.session.Close(); err != nil {
		return fmt.Errorf(
			"close Discord gateway: %w",
			err,
		)
	}

	return nil
}

func (bot *Bot) handleReady(
	session *discordgo.Session,
	event *discordgo.Ready,
) {
	slog.Info(
		"Discord bot is online",
		"user", session.State.User.Username,
		"userId", session.State.User.ID,
	)
}
