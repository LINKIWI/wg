package webgrep

import (
	"net/http"

	"wg/pkg/supercharged"
)

// Client is a webgrep API client; effectively, a single layer of abstraction above a Supercharged
// HTTP client.
type Client struct {
	sc *supercharged.HTTPClient
}

// NewClient creates a new webgrep API client for an instance hosted at a particular base URL.
func NewClient(baseURL string) (*Client, error) {
	sc, err := supercharged.NewHTTPClient(baseURL, nil)
	if err != nil {
		return nil, err
	}

	return &Client{sc}, nil
}

// Search executes a search query.
func (c *Client) Search(request *SearchQueryRequest) (*SearchQueryResponse, *supercharged.Error) {
	var resp SearchQueryResponse

	if err := c.sc.Do(http.MethodGet, "/api/search", request, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Metadata requests metadata about the webgrep instance.
func (c *Client) Metadata() (*MetadataResponse, *supercharged.Error) {
	var resp MetadataResponse

	if err := c.sc.Do(http.MethodGet, "/api/meta/info", nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
