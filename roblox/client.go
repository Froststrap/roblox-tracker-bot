package roblox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	clientSettingsCdnBaseUrl = "https://clientsettingscdn.roblox.com"
	clientSettingsBaseUrl    = "https://clientsettings.roblox.com"
	defaultChannel           = "production"
)

type Client struct {
	httpClient   *http.Client
	channelToken string
}

type clientVersionResponse struct {
	Version             string `json:"version"`
	ClientVersionUpload string `json:"clientVersionUpload"`
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (client *Client) SetChannelToken(token string) {
	client.channelToken = token
}

func normalizeChannel(channel string) string {
	if strings.EqualFold(channel, "live") ||
		strings.EqualFold(channel, "zlive") {
		return defaultChannel
	}

	return channel
}

func buildChannelPath(binaryType string, channel string) string {
	channel = normalizeChannel(channel)

	path := fmt.Sprintf(
		"/v2/client-version/%s",
		binaryType,
	)

	if !strings.EqualFold(channel, defaultChannel) {
		path += "/channel/" + url.PathEscape(channel)
	}

	return path
}

func (client *Client) getJson(
	ctx context.Context,
	path string,
	target any,
) error {
	err := client.getJsonFromBase(
		ctx,
		clientSettingsCdnBaseUrl,
		path,
		target,
	)

	if err == nil {
		return nil
	}

	if statusCode, ok := getHttpStatusCode(err); ok &&
		(statusCode == http.StatusUnauthorized ||
			statusCode == http.StatusForbidden ||
			statusCode == http.StatusNotFound) {
		return err
	}

	return client.getJsonFromBase(
		ctx,
		clientSettingsBaseUrl,
		path,
		target,
	)
}

func (client *Client) getJsonFromBase(
	ctx context.Context,
	baseUrl string,
	path string,
	target any,
) error {
	requestUrl := strings.TrimRight(baseUrl, "/") +
		"/" +
		strings.TrimLeft(path, "/")

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		requestUrl,
		nil,
	)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	if client.channelToken != "" {
		request.Header.Set(
			"Roblox-Channel-Token",
			client.channelToken,
		)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return &httpStatusError{
			statusCode: response.StatusCode,
		}
	}

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

type httpStatusError struct {
	statusCode int
}

func (err *httpStatusError) Error() string {
	return fmt.Sprintf(
		"Roblox API returned HTTP %d",
		err.statusCode,
	)
}

func getHttpStatusCode(err error) (int, bool) {
	var statusError *httpStatusError

	if !findHttpStatusError(err, &statusError) {
		return 0, false
	}

	return statusError.statusCode, true
}

func findHttpStatusError(
	err error,
	target **httpStatusError,
) bool {
	if err == nil {
		return false
	}

	if statusError, ok := err.(*httpStatusError); ok {
		*target = statusError
		return true
	}

	type wrappedError interface {
		Unwrap() error
	}

	wrapped, ok := err.(wrappedError)
	if !ok {
		return false
	}

	return findHttpStatusError(wrapped.Unwrap(), target)
}

func IsInvalidChannelStatus(statusCode int) bool {
	return statusCode == http.StatusUnauthorized ||
		statusCode == http.StatusForbidden ||
		statusCode == http.StatusNotFound
}
