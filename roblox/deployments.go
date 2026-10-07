package roblox

import (
	"context"
	"fmt"
)

type Platform string

const (
	PlatformWindows Platform = "windows"
	PlatformMacOS   Platform = "macos"
	PlatformAndroid Platform = "android"
)

type Deployment struct {
	Platform            Platform `json:"platform"`
	Channel             string   `json:"channel"`
	Version             string   `json:"version"`
	ClientVersionUpload string   `json:"clientVersionUpload"`
}

func (client *Client) GetWindowsDeployment(
	ctx context.Context,
) (Deployment, error) {
	return client.GetWindowsDeploymentForChannel(ctx, "LIVE")
}

func (client *Client) GetWindowsDeploymentForChannel(
	ctx context.Context,
	channel string,
) (Deployment, error) {
	return client.getDeployment(
		ctx,
		PlatformWindows,
		"WindowsPlayer",
		channel,
	)
}

func (client *Client) GetMacDeployment(
	ctx context.Context,
) (Deployment, error) {
	return client.GetMacDeploymentForChannel(ctx, "LIVE")
}

func (client *Client) GetMacDeploymentForChannel(
	ctx context.Context,
	channel string,
) (Deployment, error) {
	return client.getDeployment(
		ctx,
		PlatformMacOS,
		"MacPlayer",
		channel,
	)
}

func (client *Client) getDeployment(
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
