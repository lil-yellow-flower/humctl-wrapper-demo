package humanitec

import (
	"bytes"
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
	// GetApps retrieves all applications in the organization
	GetApps() ([]App, error)
	// GetApp retrieves a specific application by its ID
	GetApp(name string) (*App, error)
	// CreateApp creates a new application with the given ID and name
	// If skipEnvCreation is true, no default environment will be created
	CreateApp(id string, name string, skipEnvCreation bool) (*App, error)
	// DeleteApp deletes an application by its ID
	DeleteApp(name string) error
	// UpdateApp updates an application's name by its ID
	UpdateApp(oldName string, newName string) (*App, error)
}

// ClientFactory creates Humanitec clients
type ClientFactory interface {
	NewClient(token, org string) Client
}

// DefaultClientFactory is the default implementation of ClientFactory
type DefaultClientFactory struct{}

var (
	// defaultFactory is the default client factory
	defaultFactory ClientFactory = &DefaultClientFactory{}
)

// SetClientFactory sets the client factory (for testing only)
func SetClientFactory(factory ClientFactory) {
	defaultFactory = factory
}

// NewClient creates a new Humanitec API client
func (f *DefaultClientFactory) NewClient(token, org string) Client {
	return &humanitecClient{
		apiToken: token,
		baseURL:  "https://api.humanitec.io",
		org:      org,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewClient creates a new Humanitec API client (for backward compatibility)
func NewClient(token, org string) Client {
	return defaultFactory.NewClient(token, org)
}

// humanitecClient represents a Humanitec API client
type humanitecClient struct {
	apiToken string
	baseURL  string
	org      string
	client   *http.Client
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

// GetApp returns a specific application by its ID
func (c *humanitecClient) GetApp(name string) (*App, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/orgs/%s/apps/%s", c.baseURL, c.org, name), nil)
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

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("application not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var app App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &app, nil
}

// CreateApp creates a new application in the organization
func (c *humanitecClient) CreateApp(id string, name string, skipEnvCreation bool) (*App, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

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

// DeleteApp deletes an application by its ID
func (c *humanitecClient) DeleteApp(name string) error {
	if err := c.Validate(); err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/orgs/%s/apps/%s", c.baseURL, c.org, name), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return nil
}

// UpdateApp updates an application's name by its ID
func (c *humanitecClient) UpdateApp(oldName string, newName string) (*App, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	payload := struct {
		Name string `json:"name"`
	}{
		Name: newName,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/orgs/%s/apps/%s", c.baseURL, c.org, oldName), bytes.NewBuffer(jsonData))
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

	var app App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &app, nil
}