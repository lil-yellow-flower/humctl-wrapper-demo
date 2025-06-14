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
type mockAddClient struct {
	app   *humanitec.App
	error error
}

func (m *mockAddClient) GetApps() ([]humanitec.App, error) {
	return nil, nil
}

func (m *mockAddClient) CreateApp(name string, skipEnvCreation bool) (*humanitec.App, error) {
	if m.error != nil {
		return nil, m.error
	}
	return m.app, nil
}

// setupTestAddCommand creates a test command with the given mock client
func setupTestAddCommand(t *testing.T, mockClient *mockAddClient) *cobra.Command {
	// Create root command
	rootCmd := &cobra.Command{
		Use:   "humctl",
		Short: "Test command",
	}

	// Create add command
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add resources",
	}

	// Create app command
	appCmd := &cobra.Command{
		Use:   "app",
		Short: "Add application",
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

			// Get skip environment creation flag
			skipEnvCreation, err := cmd.Flags().GetBool(constants.SkipEnvCreationFlagName)
			if err != nil {
				return fmt.Errorf(constants.ErrInvalidSkipEnvCreation, err)
			}

			// Add application
			app, err := mockClient.CreateApp(name, skipEnvCreation)
			if err != nil {
				return fmt.Errorf(constants.ErrAddApp, err)
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

	// Add flags
	appCmd.Flags().StringP(constants.NameFlagName, constants.NameFlagShort, "", constants.NameFlagHelp)
	appCmd.MarkFlagRequired(constants.NameFlagName)
	appCmd.Flags().BoolP(constants.SkipEnvCreationFlagName, constants.SkipEnvCreationFlagShort, false, constants.SkipEnvCreationFlagHelp)
	appCmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, constants.DefaultOutputFormat, constants.OutputFlagHelp)

	// Add commands to hierarchy
	addCmd.AddCommand(appCmd)
	rootCmd.AddCommand(addCmd)

	return rootCmd
}

func TestAddApp(t *testing.T) {
	// Setup test config
	configFile := testutil.SetupTestConfig(t)
	defer os.Remove(configFile)

	// Create a mock client
	mockClient := &mockAddClient{
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
			args:           []string{"add", "app", "--name", "Test App"},
			expectedOutput: "NAME\tID\n----\t--\nTest App\ttest-app\n",
		},
		{
			name:           "json format",
			args:           []string{"add", "app", "--name", "Test App", "--output", "json"},
			expectedOutput: "{\n  \"id\": \"test-app\",\n  \"name\": \"Test App\"\n}\n",
		},
		{
			name:           "yaml format",
			args:           []string{"add", "app", "--name", "Test App", "--output", "yaml"},
			expectedOutput: "id: test-app\nname: Test App\n",
		},
		{
			name:           "successful add with skip env creation",
			args:           []string{"add", "app", "--name", "Test App", "--skip-env-creation"},
			expectedOutput: "NAME\tID\n----\t--\nTest App\ttest-app\n",
		},
		{
			name:           "successful add with short flags",
			args:           []string{"add", "app", "-n", "Test App", "-s", "-o", "json"},
			expectedOutput: "{\n  \"id\": \"test-app\",\n  \"name\": \"Test App\"\n}\n",
		},
		{
			name:        "missing name flag",
			args:        []string{"add", "app"},
			expectError: true,
		},
		{
			name:        "invalid output format",
			args:        []string{"add", "app", "--name", "Test App", "--output", "invalid"},
			expectError: true,
		},
		{
			name:        "api error",
			args:        []string{"add", "app", "--name", "Test App"},
			expectError: true,
			mockError:   fmt.Errorf("API error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set mock error if specified
			mockClient.error = tt.mockError

			// Create test command
			cmd := setupTestAddCommand(t, mockClient)

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

func TestAddCommandStructure(t *testing.T) {
	// Test add command
	assert.Equal(t, constants.AddCmdUse, addCmd.Use, "add command should have correct use")
	assert.Equal(t, constants.AddCmdShort, addCmd.Short, "add command should have correct short description")

	// Test add app command
	assert.Equal(t, constants.AddAppCmdUse, addAppCmd.Use, "add app command should have correct use")
	assert.Equal(t, constants.AddAppCmdShort, addAppCmd.Short, "add app command should have correct short description")

	// Test flags
	assert.True(t, addAppCmd.Flags().HasFlags(), "add app command should have flags")
	assert.True(t, addAppCmd.Flags().Lookup(constants.NameFlagName) != nil, "add app command should have name flag")
	assert.True(t, addAppCmd.Flags().Lookup(constants.SkipEnvCreationFlagName) != nil, "add app command should have skip-env-creation flag")
	assert.True(t, addAppCmd.Flags().Lookup(constants.OutputFlagName) != nil, "add app command should have output flag")
} 