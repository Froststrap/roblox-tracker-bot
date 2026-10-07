package roblox

import (
	"context"
	"fmt"
)

type Platform string

const (
	PlatformWindowsPlayer Platform = "WindowsPlayer"
	PlatformMacPlayer     Platform = "MacPlayer"
	PlatformWindowsStudio Platform = "WindowsStudio64"
	PlatformMacStudio     Platform = "MacStudio"
	PlatformAndroid       Platform = "Android"
)

type Deployment struct {
	Platform            Platform `json:"platform"`
	Channel             string   `json:"channel"`
	Version             string   `json:"version"`
	ClientVersionUpload string   `json:"clientVersionUpload,omitempty"`
}

func (client *Client) GetDeployment(
	ctx context.Context,
	platform Platform,
	channel string,
) (Deployment, error) {
	if channel == "" {
		channel = "LIVE"
	}

	switch platform {
	case PlatformWindowsPlayer:
		return client.getClientDeployment(
			ctx,
			platform,
			"WindowsPlayer",
			channel,
		)

	case PlatformMacPlayer:
		return client.getClientDeployment(
			ctx,
			platform,
			"MacPlayer",
			channel,
		)

	case PlatformWindowsStudio:
		return client.getClientDeployment(
			ctx,
			platform,
			"WindowsStudio64",
			channel,
		)

	case PlatformMacStudio:
		return client.getClientDeployment(
			ctx,
			platform,
			"MacStudio",
			channel,
		)

	case PlatformAndroid:
		return Deployment{}, fmt.Errorf(
			"Android deployments use the Android deployment provider",
		)

	default:
		return Deployment{}, fmt.Errorf(
			"unsupported platform: %s",
			platform,
		)
	}
}

func (client *Client) GetWindowsDeployment(
	ctx context.Context,
) (Deployment, error) {
	return client.GetDeployment(
		ctx,
		PlatformWindowsPlayer,
		"LIVE",
	)
}

func (client *Client) GetWindowsDeploymentForChannel(
	ctx context.Context,
	channel string,
) (Deployment, error) {
	return client.GetDeployment(
		ctx,
		PlatformWindowsPlayer,
		channel,
	)
}

func (client *Client) GetMacDeployment(
	ctx context.Context,
) (Deployment, error) {
	return client.GetDeployment(
		ctx,
		PlatformMacPlayer,
		"LIVE",
	)
}

func (client *Client) GetMacDeploymentForChannel(
	ctx context.Context,
	channel string,
) (Deployment, error) {
	return client.GetDeployment(
		ctx,
		PlatformMacPlayer,
		channel,
	)
}

func (client *Client) getClientDeployment(
	ctx context.Context,
	platform Platform,
	binaryType string,
	channel string,
) (Deployment, error) {
	if channel == "" {
		return Deployment{}, fmt.Errorf("channel cannot be empty")
	}

	apiChannel := normalizeChannel(channel)

	var response clientVersionResponse

	if err := client.getJson(
		ctx,
		buildChannelPath(binaryType, apiChannel),
		&response,
	); err != nil {
		return Deployment{}, fmt.Errorf(
			"get %s deployment for channel %s: %w",
			platform,
			channel,
			err,
		)
	}

	return Deployment{
		Platform:            platform,
		Channel:             channel,
		Version:             response.Version,
		ClientVersionUpload: response.ClientVersionUpload,
	}, nil
}
