package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NicolasFerreras/REPL-Pokedex/internal/model"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(url string) *Client {
	return &Client{
		BaseURL: url,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *Client) GetLocationArea(url string) (*model.LocationArea, error) {
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	var locationArea model.LocationArea
	err = json.NewDecoder(resp.Body).Decode(&locationArea)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &locationArea, nil
}
