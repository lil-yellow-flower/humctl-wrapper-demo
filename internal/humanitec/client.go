package humanitec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// App represents a Humanitec application
type App struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Client interface defines the methods that a Humanitec client must implement
type Client interface {
	// GetApps retrieves all applications in the organization
	GetApps() ([]App, error)
	// CreateApp creates a new application with the given name
	// If skipEnvCreation is true, no default environment will be created
	CreateApp(name string, skipEnvCreation bool) (*App, error)
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

// GetApps returns a list of applications
func (c *humanitecClient) GetApps() ([]App, error) {
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

// CreateApp creates a new application in the organization
func (c *humanitecClient) CreateApp(name string, skipEnvCreation bool) (*App, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Generate a URL-friendly ID from the name
	id := strings.ToLower(strings.ReplaceAll(name, " ", "-"))

	payload := struct {
		ID                    string `json:"id"`
		Name                  string `json:"name"`
		SkipEnvCreation      bool   `json:"skip_environment_creation"`
	}{
		ID:                    id,
		Name:                  name,
		SkipEnvCreation:      skipEnvCreation,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/orgs/%s/apps", c.baseURL, c.org), bytes.NewBuffer(jsonData))
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

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var app App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &app, nil
}