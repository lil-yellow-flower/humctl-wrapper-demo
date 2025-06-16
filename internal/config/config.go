package config

import (
	"fmt"
	"os"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	// HumanitecToken is the API token for authenticating with Humanitec
	HumanitecToken string `yaml:"humanitec_token"`
	// HumanitecOrg is the organization ID in Humanitec
	HumanitecOrg string `yaml:"humanitec_org"`
	// DefaultOutput is the default output format (table, json, or yaml)
	DefaultOutput string `yaml:"default_output"`
}

var (
	// Global config instance
	config Config
)

// loadConfig loads configuration from a YAML file
func loadConfig(configFile string) error {
	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("error parsing config file: %v", err)
	}

	// Set default values if not specified
	if config.DefaultOutput == "" {
		config.DefaultOutput = constants.DefaultOutputFormat
	}

	return nil
}

// GetConfig returns the current configuration
func GetConfig() Config {
	return config
}

// SetConfig sets the configuration (for testing only)
func SetConfig(c Config) {
	config = c
}

// Initialize loads the configuration
func Initialize(configFile string) error {
	return loadConfig(configFile)
} 