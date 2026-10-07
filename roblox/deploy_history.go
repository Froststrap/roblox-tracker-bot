package roblox

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
	_ "time/tzdata"
)

const (
	deployHistoryTailBytes  = 65536
	deployHistoryMaxBytes   = 4 << 20
	deployHistoryTimeLayout = "1/2/2006 3:04:05 PM"
	deployHistoryTimeZone   = "America/Los_Angeles"
)

var ErrBuildTimeNotFound = errors.New(
	"build time not found in deploy history",
)

// Matches lines such as:
// New WindowsPlayer version-hidden at 10/5/2026 3:38:32 PM, file version: ..., git hash: 0.742.0.7421053 ...Done!
var deployHistoryLinePattern = regexp.MustCompile(
	`^New (\S+) \S+ at (\d{1,2}/\d{1,2}/\d{4} \d{1,2}:\d{2}:\d{2} [AP]M),.*git hash: (\S+)`,
)

type DeployHistoryProvider struct {
	httpClient *http.Client
}

func NewDeployHistoryProvider() *DeployHistoryProvider {
	return &DeployHistoryProvider{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func deployHistorySource(platform Platform) (string, string, bool) {
	switch platform {
	case PlatformWindowsPlayer:
		return setupCdnBaseUrl + "/DeployHistory.txt", "WindowsPlayer", true
	case PlatformMacPlayer:
		return setupCdnBaseUrl + "/mac/DeployHistory.txt", "Client", true
	default:
		return "", "", false
	}
}

func (provider *DeployHistoryProvider) GetBuildTime(
	ctx context.Context,
	platform Platform,
	version string,
) (time.Time, error) {
	endpoint, entryKind, ok := deployHistorySource(platform)
	if !ok {
		return time.Time{}, ErrBuildTimeNotFound
	}

	location, err := time.LoadLocation(deployHistoryTimeZone)
	if err != nil {
		return time.Time{}, fmt.Errorf("load time zone: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return time.Time{}, fmt.Errorf("create request: %w", err)
	}

	request.Header.Set(
		"Range",
		fmt.Sprintf("bytes=-%d", deployHistoryTailBytes),
	)

	response, err := provider.httpClient.Do(request)
	if err != nil {
		return time.Time{}, fmt.Errorf("send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK &&
		response.StatusCode != http.StatusPartialContent {
		return time.Time{}, &httpStatusError{
			statusCode: response.StatusCode,
		}
	}

	body, err := io.ReadAll(
		io.LimitReader(response.Body, deployHistoryMaxBytes),
	)
	if err != nil {
		return time.Time{}, fmt.Errorf("read response: %w", err)
	}

	return findBuildTime(string(body), entryKind, version, location)
}

func findBuildTime(
	history string,
	entryKind string,
	version string,
	location *time.Location,
) (time.Time, error) {
	var found time.Time
	hasMatch := false

	for _, line := range strings.Split(history, "\n") {
		matches := deployHistoryLinePattern.FindStringSubmatch(
			strings.TrimSpace(line),
		)
		if matches == nil {
			continue
		}

		if matches[1] != entryKind || matches[3] != version {
			continue
		}

		parsed, err := time.ParseInLocation(
			deployHistoryTimeLayout,
			matches[2],
			location,
		)
		if err != nil {
			continue
		}

		// Later lines win, so a redeploy of the same version is used.
		found = parsed
		hasMatch = true
	}

	if !hasMatch {
		return time.Time{}, ErrBuildTimeNotFound
	}

	return found, nil
}
