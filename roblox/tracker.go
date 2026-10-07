package roblox

import (
	"context"
	"fmt"
	"sync"
)

type DeploymentChange struct {
	Previous *Deployment
	Current  Deployment
}

type Tracker struct {
	client          *Client
	androidProvider *AndroidDeploymentProvider

	mu    sync.RWMutex
	state map[string]Deployment
}

func NewTracker(
	client *Client,
	androidProvider *AndroidDeploymentProvider,
) *Tracker {
	return &Tracker{
		client:          client,
		androidProvider: androidProvider,
		state:           make(map[string]Deployment),
	}
}

func (tracker *Tracker) Check(
	ctx context.Context,
) ([]DeploymentChange, error) {
	deployments, err := tracker.getDeployments(ctx)
	if err != nil {
		return nil, err
	}

	changes := make([]DeploymentChange, 0)

	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	for _, deployment := range deployments {
		key := deploymentKey(deployment)

		previous, exists := tracker.state[key]

		tracker.state[key] = deployment

		if !exists {
			continue
		}

		if previous.Version == deployment.Version &&
			previous.ClientVersionUpload == deployment.ClientVersionUpload {
			continue
		}

		previousCopy := previous

		changes = append(changes, DeploymentChange{
			Previous: &previousCopy,
			Current:  deployment,
		})
	}

	return changes, nil
}

func (tracker *Tracker) getDeployments(
	ctx context.Context,
) ([]Deployment, error) {
	deployments := make([]Deployment, 0)

	platforms := []Platform{
		PlatformWindowsPlayer,
		PlatformMacPlayer,
	}

	for _, platform := range platforms {
		channels, err := tracker.client.GetPublicChannels(
			ctx,
			platform,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"get %s deployments: %w",
				platform,
				err,
			)
		}

		for _, channel := range channels {
			if channel.Deployment != nil {
				deployments = append(
					deployments,
					*channel.Deployment,
				)
			}
		}
	}

	if tracker.androidProvider != nil {
		androidDeployment, err := tracker.androidProvider.GetDeployment(ctx)
		if err != nil {
			return nil, fmt.Errorf(
				"get Android deployment: %w",
				err,
			)
		}

		deployments = append(
			deployments,
			androidDeployment,
		)
	}

	return deployments, nil
}

func deploymentKey(deployment Deployment) string {
	return string(deployment.Platform) + ":" + deployment.Channel
}
