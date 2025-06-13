package humanitec

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// App represents a Humanitec application
type App struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Client interface defines the methods that a Humanitec client must implement
type Client interface {
	ListApps() ([]App, error)
}

// humanitecClient represents a Humanitec API client
type humanitecClient struct {
	apiToken string
	baseURL  string
	org      string
	client   *http.Client
}

// NewClient creates a new Humanitec API client
func NewClient(apiToken, org string) Client {
	return &humanitecClient{
		apiToken: apiToken,
		baseURL:  "https://api.humanitec.io",
		org:      org,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Validate checks if the client is properly configured
func (c *humanitecClient) Validate() error {
	if c.apiToken == "" {
		return ErrMissingAPIToken
	}
	if c.org == "" {
		return fmt.Errorf("organization ID is required")
	}
	return nil
}

// ListApps returns a list of applications
func (c *humanitecClient) ListApps() ([]App, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/orgs/%s/apps", c.baseURL, c.org), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var apps []App
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return apps, nil
}