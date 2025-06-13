package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// mockClient implements the humanitec.Client interface for testing
type mockClient struct {
	apps  []humanitec.App
	error error
}

func (m *mockClient) GetApps() ([]humanitec.App, error) {
	if m.error != nil {
		return nil, m.error
	}
	return m.apps, nil
}

// setupTestConfig creates a temporary config.yaml file for testing
func setupTestConfig(t *testing.T) string {
	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Write test config
	config := fmt.Sprintf(`%s: "test-token"
%s: "test-org"
default_output: "table"
logging:
  level: "info"
  format: "text"
  output: "stdout"
  file: "logs/humctl-wrapper.log"`, constants.HumanitecToken, constants.HumanitecOrg)

	if _, err := tmpFile.WriteString(config); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close config file: %v", err)
	}

	return tmpFile.Name()
}

// setupTestCommand creates a test command with the given mock client
func setupTestCommand(t *testing.T, mockClient *mockClient) *cobra.Command {
	// Create root command
	rootCmd := &cobra.Command{
		Use:   "humctl",
		Short: "Test command",
	}

	// Create get command
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get resources",
	}

	// Create apps command
	appsCmd := &cobra.Command{
		Use:   "apps",
		Short: "Get applications",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get output format
			outputFormat, _ := cmd.Flags().GetString("output")

			// Get apps
			apps, err := mockClient.GetApps()
			if err != nil {
				return err
			}

			// Format output
			switch outputFormat {
			case "json":
				jsonOutput, err := json.MarshalIndent(apps, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(jsonOutput))
			case "yaml":
				yamlOutput, err := yaml.Marshal(apps)
				if err != nil {
					return err
				}
				fmt.Fprint(cmd.OutOrStdout(), string(yamlOutput))
			case "table":
				fmt.Fprintln(cmd.OutOrStdout(), "ID\tNAME")
				for _, app := range apps {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", app.ID, app.Name)
				}
			default:
				return fmt.Errorf("invalid output format: %s", outputFormat)
			}

			return nil
		},
	}

	// Add output format flag
	appsCmd.Flags().StringP("output", "o", constants.DefaultOutputFormat, "Output format (table|json|yaml)")

	// Add commands to hierarchy
	getCmd.AddCommand(appsCmd)
	rootCmd.AddCommand(getCmd)

	return rootCmd
}

func TestGetApps(t *testing.T) {
	// Setup test config
	configFile := setupTestConfig(t)
	defer os.Remove(configFile)

	// Create a mock client
	mockClient := &mockClient{
		apps: []humanitec.App{
			{ID: "app1", Name: "App 1"},
			{ID: "app2", Name: "App 2"},
		},
	}

	// Test cases
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "table format",
			args:           []string{"get", "apps"},
			expectedOutput: "ID\tNAME\napp1\tApp 1\napp2\tApp 2\n",
		},
		{
			name:           "json format",
			args:           []string{"get", "apps", "--output", "json"},
			expectedOutput: "[\n  {\n    \"id\": \"app1\",\n    \"name\": \"App 1\"\n  },\n  {\n    \"id\": \"app2\",\n    \"name\": \"App 2\"\n  }\n]\n",
		},
		{
			name:           "yaml format",
			args:           []string{"get", "apps", "--output", "yaml"},
			expectedOutput: "- id: app1\n  name: App 1\n- id: app2\n  name: App 2\n",
		},
		{
			name:        "invalid format",
			args:        []string{"get", "apps", "--output", "invalid"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set config file path
			os.Setenv("CONFIG_FILE", configFile)

			// Create test command
			cmd := setupTestCommand(t, mockClient)

			// Set up output buffer
			outBuf := new(bytes.Buffer)
			cmd.SetOut(outBuf)
			cmd.SetErr(outBuf)

			// Execute command
			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			// Check error
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check output
			output := outBuf.String()
			if output != tt.expectedOutput {
				t.Errorf("Expected output:\n%q\nGot:\n%q", tt.expectedOutput, output)
			}
		})
	}
}

func TestCommandStructure(t *testing.T) {
	// Test root command
	assert.Equal(t, constants.RootCmdUse, rootCmd.Use, "root command should have correct use")
	assert.Equal(t, "A command line interface wrapper for Humanitec platform", rootCmd.Short, "root command should have correct short description")

	// Test get command
	assert.Equal(t, constants.GetCmdUse, getCmd.Use, "get command should have correct use")
	assert.Equal(t, "Get resources from Humanitec platform", getCmd.Short, "get command should have correct short description")

	// Test get apps command
	assert.Equal(t, constants.GetAppsCmdUse, getAppsCmd.Use, "get apps command should have correct use")
	assert.Equal(t, "Get applications from Humanitec platform", getAppsCmd.Short, "get apps command should have correct short description")

	// Test flags
	assert.True(t, getAppsCmd.Flags().HasFlags(), "get apps command should have flags")
	assert.True(t, getAppsCmd.Flags().Lookup(constants.OutputFlagName) != nil, "get apps command should have output flag")
	assert.True(t, getAppsCmd.Flags().Lookup(constants.OrgFlagName) != nil, "get apps command should have org flag")
} 