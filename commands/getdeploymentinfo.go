package commands

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Froststrap/roblox-tracker-bot/roblox"
	"github.com/bwmarrin/discordgo"
)

const (
	assetBaseUrl = "https://raw.githubusercontent.com/" +
		"Froststrap/roblox-tracker-bot/main/assets/"
	playerIconUrl = assetBaseUrl + "robloxPlayer.png"
	studioIconUrl = assetBaseUrl + "robloxStudio.png"

	defaultEmbedColor = 0x335FFF
	androidEmbedColor = 0x3DDC84
)

func GetDeploymentInfoCommand() *Command {
	deploymentService := roblox.NewDeploymentService()

	return &Command{
		Name:        "getdeploymentinfo",
		Description: "Get the current Roblox deployment information",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "platform",
				Description: "Roblox platform",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Windows Player",
						Value: string(roblox.PlatformWindowsPlayer),
					},
					{
						Name:  "Mac Player",
						Value: string(roblox.PlatformMacPlayer),
					},
					{
						Name:  "Windows Studio",
						Value: string(roblox.PlatformWindowsStudio),
					},
					{
						Name:  "Mac Studio",
						Value: string(roblox.PlatformMacStudio),
					},
					{
						Name:  "Android",
						Value: string(roblox.PlatformAndroid),
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "channel",
				Description: "Roblox channel",
				Required:    false,
			},
		},
		ExecuteInteraction: func(commandContext *CommandContext) error {
			return executeGetDeploymentInfo(
				commandContext,
				deploymentService,
			)
		},
	}
}

func executeGetDeploymentInfo(
	commandContext *CommandContext,
	deploymentService *roblox.DeploymentService,
) error {
	options := commandContext.Interaction.ApplicationCommandData().Options

	platform := roblox.Platform(
		getOptionString(options, "platform"),
	)

	channel := getOptionString(options, "channel")

	if strings.TrimSpace(channel) == "" {
		channel = "LIVE"
	}

	// Fetching deployment details can take longer than Discord
	err := commandContext.Session.InteractionRespond(
		commandContext.Interaction.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"defer deployment info response: %w",
			err,
		)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancel()

	details, err := deploymentService.GetDeploymentDetails(
		ctx,
		platform,
		channel,
	)
	if err != nil {
		slog.Error(
			"failed to get deployment info",
			"platform", platform,
			"channel", channel,
			"error", err,
		)

		return editDeploymentResponse(
			commandContext,
			&discordgo.WebhookEdit{
				Content: stringPointer(
					"Could not get deployment information. " +
						"Check the platform and channel and try again.",
				),
			},
		)
	}

	embeds := []*discordgo.MessageEmbed{
		buildDeploymentEmbed(platform, details),
	}
	components := buildDownloadComponents(details.DownloadLinks)

	return editDeploymentResponse(
		commandContext,
		&discordgo.WebhookEdit{
			Embeds:     &embeds,
			Components: &components,
		},
	)
}

func getPlatformIconUrl(platform roblox.Platform) string {
	switch platform {
	case roblox.PlatformWindowsStudio, roblox.PlatformMacStudio:
		return studioIconUrl
	default:
		return playerIconUrl
	}
}

func getPlatformEmbedColor(platform roblox.Platform) int {
	if platform == roblox.PlatformAndroid {
		return androidEmbedColor
	}

	return defaultEmbedColor
}

func buildDeploymentEmbed(
	platform roblox.Platform,
	details roblox.DeploymentDetails,
) *discordgo.MessageEmbed {
	deployment := details.Deployment

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Current version for %s", platform),
		URL:   details.ClientSettingsUrl,
		Description: fmt.Sprintf(
			"The `%s` channel's current deployment for `%s`.",
			deployment.Channel,
			platform,
		),
		Color: getPlatformEmbedColor(platform),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: getPlatformIconUrl(platform),
		},
	}

	if details.BuiltAt != nil {
		unix := details.BuiltAt.Unix()

		embed.Fields = append(
			embed.Fields,
			&discordgo.MessageEmbedField{
				Name: "Built on:",
				Value: fmt.Sprintf(
					"<t:%d:F> (<t:%d:R>)\nFound in deployment history",
					unix,
					unix,
				),
			},
		)
	}

	embed.Fields = append(
		embed.Fields,
		&discordgo.MessageEmbedField{
			Name:   "Version:",
			Value:  "```" + deployment.Version + "```",
			Inline: true,
		},
	)

	if deployment.ClientVersionUpload != "" {
		embed.Fields = append(
			embed.Fields,
			&discordgo.MessageEmbedField{
				Name:   "GUID/Hash:",
				Value:  "```" + deployment.ClientVersionUpload + "```",
				Inline: true,
			},
		)
	}

	return embed
}

func buildDownloadComponents(
	links []roblox.DownloadLink,
) []discordgo.MessageComponent {
	if len(links) == 0 {
		return []discordgo.MessageComponent{}
	}

	buttons := make([]discordgo.MessageComponent, 0, len(links))

	for _, link := range links {
		buttons = append(
			buttons,
			discordgo.Button{
				Label: fmt.Sprintf("Download (%s)", link.Name),
				Style: discordgo.LinkButton,
				URL:   link.Url,
			},
		)
	}

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: buttons,
		},
	}
}

func editDeploymentResponse(
	commandContext *CommandContext,
	edit *discordgo.WebhookEdit,
) error {
	_, err := commandContext.Session.InteractionResponseEdit(
		commandContext.Interaction.Interaction,
		edit,
	)
	if err != nil {
		return fmt.Errorf(
			"edit deployment info response: %w",
			err,
		)
	}

	return nil
}

func stringPointer(value string) *string {
	return &value
}

func getOptionString(
	options []*discordgo.ApplicationCommandInteractionDataOption,
	name string,
) string {
	for _, option := range options {
		if option.Name == name {
			return option.StringValue()
		}
	}

	return ""
}
