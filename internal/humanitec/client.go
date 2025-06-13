package humanitec

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// App represents a Humanitec application
type App struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Client represents a Humanitec API client
type Client struct {
	apiToken string
	baseURL  string
	org      string
	env      string
	client   *http.Client
}

// NewClient creates a new Humanitec API client
func NewClient() *Client {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found")
	}

	return &Client{
		apiToken: os.Getenv("HUMANITEC_TOKEN"),
		baseURL:  "https://api.humanitec.io",
		org:      os.Getenv("HUMANITEC_ORG"),
		env:      os.Getenv("HUMANITEC_ENV"),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Validate checks if the client is properly configured
func (c *Client) Validate() error {
	if c.apiToken == "" {
		return ErrMissingAPIToken
	}
	if c.org == "" {
		return fmt.Errorf("HUMANITEC_ORG environment variable is not set")
	}
	return nil
}

// ListApps returns a list of applications
func (c *Client) ListApps() ([]App, error) {
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