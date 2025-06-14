package commands

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/output"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/testutil"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// mockClient implements the humanitec.Client interface for testing
type mockClient struct {
	apps  []humanitec.App
	app   *humanitec.App
	error error
}

func (m *mockClient) GetApps() ([]humanitec.App, error) {
	if m.error != nil {
		return nil, m.error
	}
	return m.apps, nil
}

func (m *mockClient) CreateApp(name string, skipEnvCreation bool) (*humanitec.App, error) {
	return nil, nil
}

func (m *mockClient) GetApp(name string) (*humanitec.App, error) {
	if m.error != nil {
		return nil, m.error
	}
	return m.app, nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get format from flag
			formatStr, err := cmd.Flags().GetString(constants.OutputFlagName)
			if err != nil {
				return fmt.Errorf(constants.ErrInvalidOutputFormat, err)
			}

			// Validate format
			format, err := output.ValidateFormat(formatStr)
			if err != nil {
				return fmt.Errorf(constants.ErrInvalidOutputFormat, err)
			}

			// Get apps
			apps, err := mockClient.GetApps()
			if err != nil {
				return err
			}

			// Format and print output
			formatted, err := output.FormatApps(apps, format)
			if err != nil {
				return fmt.Errorf(constants.ErrFormatOutput, err)
			}
			fmt.Fprint(cmd.OutOrStdout(), formatted)

			return nil
		},
	}

	// Create get app command
	getAppCmd := &cobra.Command{
		Use:   constants.GetAppCmdUse,
		Short: constants.GetAppCmdShort,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get format from flag
			formatStr, err := cmd.Flags().GetString(constants.OutputFlagName)
			if err != nil {
				return fmt.Errorf(constants.ErrInvalidOutputFormat, err)
			}

			// Validate format
			format, err := output.ValidateFormat(formatStr)
			if err != nil {
				return fmt.Errorf(constants.ErrInvalidOutputFormat, err)
			}

			// Get name from flag
			name, err := cmd.Flags().GetString(constants.NameFlagName)
			if err != nil {
				return fmt.Errorf(constants.ErrInvalidName, err)
			}

			// Get application
			app, err := mockClient.GetApp(name)
			if err != nil {
				return fmt.Errorf(constants.ErrGetApp, err)
			}

			// Format and print output
			formatted, err := output.FormatApp(app, format)
			if err != nil {
				return fmt.Errorf(constants.ErrFormatOutput, err)
			}
			fmt.Fprint(cmd.OutOrStdout(), formatted)

			return nil
		},
	}

	// Add flags to get app command
	getAppCmd.Flags().StringP(constants.NameFlagName, constants.NameFlagShort, "", constants.NameFlagHelp)
	getAppCmd.MarkFlagRequired(constants.NameFlagName)
	getAppCmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, constants.DefaultOutputFormat, constants.OutputFlagHelp)

	// Add output format flag to get command
	getCmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, constants.DefaultOutputFormat, constants.OutputFlagHelp)

	// Add commands to hierarchy
	getCmd.AddCommand(getAppCmd)
	rootCmd.AddCommand(getCmd)

	return rootCmd
}

func TestGetApps(t *testing.T) {
	// Setup test config
	configFile := testutil.SetupTestConfig(t)
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
			args:           []string{"get"},
			expectedOutput: "NAME\tID\n----\t--\nApp 1\tapp1\nApp 2\tapp2\n",
		},
		{
			name:           "json format",
			args:           []string{"get", "--output", "json"},
			expectedOutput: "[\n  {\n    \"id\": \"app1\",\n    \"name\": \"App 1\"\n  },\n  {\n    \"id\": \"app2\",\n    \"name\": \"App 2\"\n  }\n]\n",
		},
		{
			name:           "yaml format",
			args:           []string{"get", "--output", "yaml"},
			expectedOutput: "- id: app1\n  name: App 1\n- id: app2\n  name: App 2\n",
		},
		{
			name:        "invalid format",
			args:        []string{"get", "--output", "invalid"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestGetApp(t *testing.T) {
	// Setup test config
	configFile := testutil.SetupTestConfig(t)
	defer os.Remove(configFile)

	// Create a mock client
	mockClient := &mockClient{
		app: &humanitec.App{
			ID:   "test-app",
			Name: "Test App",
		},
	}

	// Test cases
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
		mockError      error
	}{
		{
			name:           "table format",
			args:           []string{"get", "app", "--name", "test-app"},
			expectedOutput: "NAME\tID\n----\t--\nTest App\ttest-app\n",
		},
		{
			name:           "json format",
			args:           []string{"get", "app", "--name", "test-app", "--output", "json"},
			expectedOutput: "{\n  \"id\": \"test-app\",\n  \"name\": \"Test App\"\n}\n",
		},
		{
			name:           "yaml format",
			args:           []string{"get", "app", "--name", "test-app", "--output", "yaml"},
			expectedOutput: "id: test-app\nname: Test App\n",
		},
		{
			name:        "missing name flag",
			args:        []string{"get", "app"},
			expectError: true,
		},
		{
			name:        "invalid output format",
			args:        []string{"get", "app", "--name", "test-app", "--output", "invalid"},
			expectError: true,
		},
		{
			name:        "api error",
			args:        []string{"get", "app", "--name", "test-app"},
			expectError: true,
			mockError:   fmt.Errorf("API error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test command
			cmd := setupTestCommand(t, mockClient)
			cmd.SetArgs(tt.args)

			// Set up output buffer
			outBuf := new(bytes.Buffer)
			cmd.SetOut(outBuf)
			cmd.SetErr(outBuf)

			// Set mock error if specified
			if tt.mockError != nil {
				mockClient.error = tt.mockError
			} else {
				mockClient.error = nil
			}

			// Execute the command
			err := cmd.Execute()

			// Check error expectations
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
	assert.Equal(t, constants.RootCmdShort, rootCmd.Short, "root command should have correct short description")

	// Test get command
	assert.Equal(t, constants.GetCmdUse, getCmd.Use, "get command should have correct use")
	assert.Equal(t, constants.GetCmdShort, getCmd.Short, "get command should have correct short description")

	// Test flags
	assert.True(t, getCmd.Flags().HasFlags(), "get command should have flags")
	assert.True(t, getCmd.Flags().Lookup(constants.OutputFlagName) != nil, "get command should have output flag")
} 