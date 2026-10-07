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
	Platform            Platform
	Channel             string
	Version             string
	ClientVersionUpload string
}

func (client *Client) GetWindowsDeployment(ctx context.Context) (Deployment, error) {
	return client.GetWindowsDeploymentForChannel(ctx, "LIVE")
}

func (client *Client) GetWindowsDeploymentForChannel(
	ctx context.Context,
	channel string,
) (Deployment, error) {
	var response struct {
		Version             string `json:"version"`
		ClientVersionUpload string `json:"clientVersionUpload"`
	}

	path := fmt.Sprintf(
		"/v2/client-version/WindowsPlayer/channel/%s",
		channel,
	)

	if err := client.getJson(ctx, path, &response); err != nil {
		return Deployment{}, fmt.Errorf(
			"get Windows deployment for channel %s: %w",
			channel,
			err,
		)
	}

	return Deployment{
		Platform:            PlatformWindows,
		Channel:             channel,
		Version:             response.Version,
		ClientVersionUpload: response.ClientVersionUpload,
	}, nil
}

func (client *Client) GetMacDeployment(ctx context.Context) (Deployment, error) {
	return client.GetMacDeploymentForChannel(ctx, "LIVE")
}

func (client *Client) GetMacDeploymentForChannel(
	ctx context.Context,
	channel string,
) (Deployment, error) {
	var response struct {
		Version             string `json:"version"`
		ClientVersionUpload string `json:"clientVersionUpload"`
	}

	path := fmt.Sprintf(
		"/v2/client-version/MacPlayer/channel/%s",
		channel,
	)

	if err := client.getJson(ctx, path, &response); err != nil {
		return Deployment{}, fmt.Errorf(
			"get macOS deployment for channel %s: %w",
			channel,
			err,
		)
	}

	return Deployment{
		Platform:            PlatformMacOS,
		Channel:             channel,
		Version:             response.Version,
		ClientVersionUpload: response.ClientVersionUpload,
	}, nil
}
