package roblox

import (
	"context"
	"testing"
)

func TestGetWindowsDeployment(t *testing.T) {
	client := NewClient()

	deployment, err := client.GetWindowsDeployment(context.Background())
	if err != nil {
		t.Fatalf("failed to get Windows deployment: %v", err)
	}

	if deployment.Platform != PlatformWindows {
		t.Fatalf("expected platform %q, got %q", PlatformWindows, deployment.Platform)
	}

	if deployment.Channel != "LIVE" {
		t.Fatalf("expected channel LIVE, got %q", deployment.Channel)
	}

	if deployment.Version == "" {
		t.Fatal("Windows deployment returned an empty version")
	}

	if deployment.ClientVersionUpload == "" {
		t.Fatal("Windows deployment returned an empty clientVersionUpload")
	}

	t.Logf(
		"Windows deployment: %s / %s / %s",
		deployment.Channel,
		deployment.Version,
		deployment.ClientVersionUpload,
	)
}

func TestGetMacDeployment(t *testing.T) {
	client := NewClient()

	deployment, err := client.GetMacDeployment(context.Background())
	if err != nil {
		t.Fatalf("failed to get macOS deployment: %v", err)
	}

	if deployment.Platform != PlatformMacOS {
		t.Fatalf("expected platform %q, got %q", PlatformMacOS, deployment.Platform)
	}

	if deployment.Channel != "LIVE" {
		t.Fatalf("expected channel LIVE, got %q", deployment.Channel)
	}

	if deployment.Version == "" {
		t.Fatal("macOS deployment returned an empty version")
	}

	if deployment.ClientVersionUpload == "" {
		t.Fatal("macOS deployment returned an empty clientVersionUpload")
	}

	t.Logf(
		"macOS deployment: %s / %s / %s",
		deployment.Channel,
		deployment.Version,
		deployment.ClientVersionUpload,
	)
}

func TestGetPublicChannels(t *testing.T) {
	client := NewClient()

	channels, err := client.GetPublicChannels(
		context.Background(),
		PlatformWindows,
	)
	if err != nil {
		t.Fatalf("failed to get public channels: %v", err)
	}

	if len(channels) == 0 {
		t.Fatal("no public channels found")
	}

	for _, channel := range channels {
		if channel.Name == "" {
			t.Fatal("public channel returned an empty name")
		}

		if channel.Deployment == nil {
			t.Fatalf(
				"public channel %s returned no deployment",
				channel.Name,
			)
		}

		t.Logf(
			"public channel: %s -> %s (%s)",
			channel.Name,
			channel.Deployment.Version,
			channel.Deployment.ClientVersionUpload,
		)
	}
}
