package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func PingCommand() *Command {
	return &Command{
		Name:               "ping",
		Description:        "Check whether RbxWatch is responding",
		ExecuteInteraction: executePing,
	}
}

func executePing(commandContext *CommandContext) error {
	latency := commandContext.Session.HeartbeatLatency()

	err := commandContext.Session.InteractionRespond(
		commandContext.Interaction.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(
					"Pong! `%dms`",
					latency.Milliseconds(),
				),
			},
		},
	)
	if err != nil {
		return fmt.Errorf(
			"respond to ping interaction: %w",
			err,
		)
	}

	return nil
}
