package roblox

import (
	"context"
	"fmt"
)

type Channel struct {
	Name       string      `json:"name"`
	Public     bool        `json:"public"`
	Deployment *Deployment `json:"deployment,omitempty"`
}

var knownChannels = []string{
	"LIVE",
	"zCanary",
	"zIntegration",
}

func GetKnownChannels() []string {
	channels := make([]string, len(knownChannels))
	copy(channels, knownChannels)

	return channels
}

func (client *Client) GetPublicChannels(
	ctx context.Context,
	platform Platform,
) ([]Channel, error) {
	binaryType, err := getBinaryType(platform)
	if err != nil {
		return nil, err
	}

	channels := make([]Channel, 0, len(knownChannels))

	for _, channelName := range knownChannels {
		deployment, err := client.getDeployment(
			ctx,
			platform,
			binaryType,
			channelName,
		)

		if err != nil {
			if isInvalidChannelError(err) {
				continue
			}

			return nil, fmt.Errorf(
				"check channel %s: %w",
				channelName,
				err,
			)
		}

		channels = append(channels, Channel{
			Name:       channelName,
			Public:     true,
			Deployment: &deployment,
		})
	}

	return channels, nil
}

func getBinaryType(platform Platform) (string, error) {
	switch platform {
	case PlatformWindows:
		return "WindowsPlayer", nil
	case PlatformMacOS:
		return "MacPlayer", nil
	default:
		return "", fmt.Errorf(
			"unsupported platform: %s",
			platform,
		)
	}
}

func isInvalidChannelError(err error) bool {
	statusCode, ok := getHttpStatusCode(err)

	if !ok {
		return false
	}

	return IsInvalidChannelStatus(statusCode)
}
