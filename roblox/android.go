package roblox

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const androidGooglePlayUrl = "https://play.google.com/store/apps/details?id=com.roblox.client"

type AndroidDeploymentProvider struct {
	httpClient *http.Client
}

func NewAndroidDeploymentProvider() *AndroidDeploymentProvider {
	return &AndroidDeploymentProvider{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (provider *AndroidDeploymentProvider) GetDeployment(
	ctx context.Context,
) (Deployment, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		androidGooglePlayUrl,
		nil,
	)
	if err != nil {
		return Deployment{}, fmt.Errorf(
			"create Android deployment request: %w",
			err,
		)
	}

	request.Header.Set(
		"User-Agent",
		"Mozilla/5.0",
	)

	response, err := provider.httpClient.Do(request)
	if err != nil {
		return Deployment{}, fmt.Errorf(
			"request Android deployment: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Deployment{}, fmt.Errorf(
			"Google Play returned HTTP %d",
			response.StatusCode,
		)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return Deployment{}, fmt.Errorf(
			"parse Google Play response: %w",
			err,
		)
	}

	version, err := findAndroidVersion(document)
	if err != nil {
		return Deployment{}, err
	}

	return Deployment{
		Platform: PlatformAndroid,
		Channel:  "LIVE",
		Version:  version,
	}, nil
}

func findAndroidVersion(document *goquery.Document) (string, error) {
	versionPattern := regexp.MustCompile(`\b2\.\d+\.\d+\b`)

	var version string

	document.Find("body").EachWithBreak(func(_ int, selection *goquery.Selection) bool {
		matches := versionPattern.FindAllString(selection.Text(), -1)

		if len(matches) == 0 {
			return true
		}

		version = matches[len(matches)-1]
		return false
	})

	if version == "" {
		return "", fmt.Errorf(
			"could not find Roblox Android version on Google Play",
		)
	}

	return version, nil
}
