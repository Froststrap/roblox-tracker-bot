package roblox

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type DeploymentDetails struct {
	Deployment        Deployment
	ClientSettingsUrl string
	BuiltAt           *time.Time
	DownloadLinks     []DownloadLink
}

type DeploymentService struct {
	client          *Client
	androidProvider *AndroidDeploymentProvider
	historyProvider *DeployHistoryProvider
}

func NewDeploymentService() *DeploymentService {
	return &DeploymentService{
		client:          NewClient(),
		androidProvider: NewAndroidDeploymentProvider(),
		historyProvider: NewDeployHistoryProvider(),
	}
}

func (service *DeploymentService) GetDeployment(
	ctx context.Context,
	platform Platform,
	channel string,
) (Deployment, error) {
	if platform == PlatformAndroid {
		deployment, err := service.androidProvider.GetDeployment(ctx)
		if err != nil {
			return Deployment{}, fmt.Errorf(
				"get Android deployment: %w",
				err,
			)
		}

		return deployment, nil
	}

	deployment, err := service.client.GetDeployment(
		ctx,
		platform,
		channel,
	)
	if err != nil {
		return Deployment{}, fmt.Errorf(
			"get %s deployment: %w",
			platform,
			err,
		)
	}

	return deployment, nil
}

func (service *DeploymentService) GetDeploymentDetails(
	ctx context.Context,
	platform Platform,
	channel string,
) (DeploymentDetails, error) {
	deployment, err := service.GetDeployment(ctx, platform, channel)
	if err != nil {
		return DeploymentDetails{}, err
	}

	details := DeploymentDetails{
		Deployment: deployment,
		ClientSettingsUrl: ClientSettingsUrl(
			platform,
			deployment.Channel,
		),
	}

	if !isProductionChannel(deployment.Channel) {
		return details, nil
	}

	details.DownloadLinks = DownloadLinks(platform, deployment)

	builtAt, err := service.historyProvider.GetBuildTime(
		ctx,
		platform,
		deployment.Version,
	)
	if err == nil {
		details.BuiltAt = &builtAt
	} else if !errors.Is(err, ErrBuildTimeNotFound) {
		slog.Warn(
			"could not read deploy history",
			"platform", platform,
			"error", err,
		)
	}

	return details, nil
}
