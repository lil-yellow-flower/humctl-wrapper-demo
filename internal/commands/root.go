package commands

import (
	"fmt"
	"os"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/spf13/cobra"
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
	version = "0.3.0"
	commit  = "dev"
	date    = "unknown"

	// Global config instance
	config Config

	rootCmd = &cobra.Command{
		Use:     constants.RootCmdUse,
		Short:   "A command line interface wrapper for Humanitec platform",
		Long:    `A command line interface wrapper for Humanitec platform that provides basic CRUD operations for managing resources.`,
		Version: fmt.Sprintf("%s", version),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Load config file
			if err := loadConfig("config.yaml"); err != nil {
				return fmt.Errorf("error loading config: %v", err)
			}

			// Verify required configuration values are set
			if config.HumanitecToken == "" {
				return fmt.Errorf(constants.ErrMissingToken)
			}
			if config.HumanitecOrg == "" {
				return fmt.Errorf(constants.ErrMissingOrg)
			}

			return nil
		},
	}
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

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	// Add get command to root
	rootCmd.AddCommand(getCmd)
	return rootCmd.Execute()
}

func init() {
	// Add version flag
	rootCmd.PersistentFlags().BoolP(constants.VersionFlagName, constants.VersionFlagShort, false, "Print the version number")
} 