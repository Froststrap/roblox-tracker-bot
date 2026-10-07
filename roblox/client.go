package roblox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseUrl = "https://clientsettingscdn.roblox.com"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (client *Client) getJson(
	ctx context.Context,
	path string,
	target any,
) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		baseUrl+path,
		nil,
	)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Roblox API returned HTTP %d", response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
