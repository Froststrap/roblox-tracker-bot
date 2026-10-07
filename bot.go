package main

import (
	"fmt"
	"log/slog"

	"github.com/Froststrap/roblox-tracker-bot/commands"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session  *discordgo.Session
	config   *Config
	registry *commands.Registry
}

func newBot(config *Config) (*Bot, error) {
	session, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		return nil, fmt.Errorf(
			"create Discord session: %w",
			err,
		)
	}

	registry := commands.NewRegistry()

	// Register commands here
	registry.Register(commands.PingCommand())

	bot := &Bot{
		session:  session,
		config:   config,
		registry: registry,
	}

	session.AddHandler(bot.handleReady)
	session.AddHandler(bot.handleInteractionCreate)

	session.Identify.Intents = discordgo.IntentsGuilds

	return bot, nil
}

func (bot *Bot) start() error {
	if err := bot.session.Open(); err != nil {
		return fmt.Errorf(
			"open Discord gateway: %w",
			err,
		)
	}

	if err := bot.registerApplicationCommands(); err != nil {
		if closeError := bot.session.Close(); closeError != nil {
			slog.Error(
				"failed to close Discord session after command registration failure",
				"error", closeError,
			)
		}

		return err
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

func (bot *Bot) registerApplicationCommands() error {
	applicationId := bot.session.State.User.ID
	applicationCommands := bot.registry.ApplicationCommands()

	_, err := bot.session.ApplicationCommandBulkOverwrite(
		applicationId,
		bot.config.DiscordGuildId,
		applicationCommands,
	)
	if err != nil {
		return fmt.Errorf(
			"register application commands: %w",
			err,
		)
	}

	slog.Info(
		"synchronised application commands",
		"count", len(applicationCommands),
		"guildId", bot.config.DiscordGuildId,
	)

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

func (bot *Bot) handleInteractionCreate(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandName := interaction.ApplicationCommandData().Name
	command := bot.registry.Find(commandName)

	if command == nil || command.ExecuteInteraction == nil {
		return
	}

	userId := ""

	if interaction.Member != nil &&
		interaction.Member.User != nil {
		userId = interaction.Member.User.ID
	} else if interaction.User != nil {
		userId = interaction.User.ID
	}

	slog.Info(
		"received application command",
		"command", commandName,
		"userId", userId,
		"guildId", interaction.GuildID,
	)

	commandContext := &commands.CommandContext{
		Session:     session,
		Interaction: interaction,
	}

	if err := command.ExecuteInteraction(commandContext); err != nil {
		slog.Error(
			"application command failed",
			"command", commandName,
			"error", err,
		)

		_ = session.InteractionRespond(
			interaction.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "An error occurred while executing that command.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			},
		)
	}
}
