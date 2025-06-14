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

// mockUpdateClient implements the humanitec.Client interface for testing
type mockUpdateClient struct {
	error error
}

func (m *mockUpdateClient) GetApps() ([]humanitec.App, error) {
	return nil, nil
}

func (m *mockUpdateClient) CreateApp(name string, skipEnvCreation bool) (*humanitec.App, error) {
	return nil, nil
}

func (m *mockUpdateClient) DeleteApp(name string) error {
	return nil
}

func (m *mockUpdateClient) UpdateApp(name string, newName string) error {
	return m.error
}

// setupTestUpdateCommand creates a test command with the given mock client
func setupTestUpdateCommand(t *testing.T, mockClient *mockUpdateClient) *cobra.Command {
	// Create root command
	rootCmd := &cobra.Command{
		Use:   "humctl",
		Short: "Test command",
	}

	// Create update command
	updateCmd := &cobra.Command{
		Use:   constants.UpdateCmdUse,
		Short: constants.UpdateCmdShort,
	}

	// Create app command
	appCmd := &cobra.Command{
		Use:   constants.UpdateAppCmdUse,
		Short: constants.UpdateAppCmdShort,
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

			// Get new name from flag
			newName, err := cmd.Flags().GetString(constants.NewNameFlagName)
			if err != nil {
				return fmt.Errorf(constants.ErrInvalidNewName, err)
			}

			// Update application
			err = mockClient.UpdateApp(name, newName)
			if err != nil {
				return fmt.Errorf(constants.ErrUpdateApp, err)
			}

			// Format and print output
			formatted, err := output.FormatMessage("Application successfully updated", format)
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

	appCmd.Flags().StringP(constants.NewNameFlagName, constants.NewNameFlagShort, "", constants.NewNameFlagHelp)
	appCmd.MarkFlagRequired(constants.NewNameFlagName)

	appCmd.Flags().StringP(constants.OrgFlagName, constants.OrgFlagShort, "", fmt.Sprintf("Humanitec organization ID (defaults to %s environment variable)", constants.HumanitecOrg))
	appCmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, constants.DefaultOutputFormat, constants.OutputFlagHelp)

	// Add commands to hierarchy
	updateCmd.AddCommand(appCmd)
	rootCmd.AddCommand(updateCmd)

	return rootCmd
}

func TestUpdateApp(t *testing.T) {
	// Setup test config
	configFile := testutil.SetupTestConfig(t)
	defer os.Remove(configFile)

	// Create a mock client
	mockClient := &mockUpdateClient{}

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
			args:           []string{"update", "app", "--name", "test-app", "--new-name", "Updated Test App"},
			expectedOutput: constants.SuccessAppUpdated + "\n",
		},
		{
			name:           "json format",
			args:           []string{"update", "app", "--name", "test-app", "--new-name", "Updated Test App", "--output", "json"},
			expectedOutput: "{\n  \"message\": \"" + constants.SuccessAppUpdated + "\"\n}\n",
		},
		{
			name:           "yaml format",
			args:           []string{"update", "app", "--name", "test-app", "--new-name", "Updated Test App", "--output", "yaml"},
			expectedOutput: "message: " + constants.SuccessAppUpdated + "\n",
		},
		{
			name:           "successful update with short flags",
			args:           []string{"update", "app", "-n", "test-app", "-m", "Updated Test App", "-o", "json"},
			expectedOutput: "{\n  \"message\": \"" + constants.SuccessAppUpdated + "\"\n}\n",
		},
		{
			name:           "successful update with organization flag",
			args:           []string{"update", "app", "--name", "test-app", "--new-name", "Updated Test App", "--org", "test-org"},
			expectedOutput: constants.SuccessAppUpdated + "\n",
		},
		{
			name:           "successful update with organization short flag",
			args:           []string{"update", "app", "-n", "test-app", "-m", "Updated Test App", "-g", "test-org"},
			expectedOutput: constants.SuccessAppUpdated + "\n",
		},
		{
			name:        "missing name flag",
			args:        []string{"update", "app", "--new-name", "Updated Test App"},
			expectError: true,
		},
		{
			name:        "missing new-name flag",
			args:        []string{"update", "app", "--name", "test-app"},
			expectError: true,
		},
		{
			name:        "invalid output format",
			args:        []string{"update", "app", "--name", "test-app", "--new-name", "Updated Test App", "--output", "invalid"},
			expectError: true,
		},
		{
			name:        "api error",
			args:        []string{"update", "app", "--name", "test-app", "--new-name", "Updated Test App"},
			expectError: true,
			mockError:   fmt.Errorf(constants.ErrAPIError),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set mock error if specified
			mockClient.error = tt.mockError

			// Create test command
			cmd := setupTestUpdateCommand(t, mockClient)

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

func TestUpdateCommandStructure(t *testing.T) {
	// Test update command
	assert.Equal(t, constants.UpdateCmdUse, updateCmd.Use, "update command should have correct use")
	assert.Equal(t, constants.UpdateCmdShort, updateCmd.Short, "update command should have correct short description")

	// Test update app command
	assert.Equal(t, constants.UpdateAppCmdUse, updateAppCmd.Use, "update app command should have correct use")
	assert.Equal(t, constants.UpdateAppCmdShort, updateAppCmd.Short, "update app command should have correct short description")

	// Test flags
	assert.True(t, updateAppCmd.Flags().HasFlags(), "update app command should have flags")
	assert.True(t, updateAppCmd.Flags().Lookup(constants.NameFlagName) != nil, "update app command should have name flag")
	assert.True(t, updateAppCmd.Flags().Lookup(constants.NewNameFlagName) != nil, "update app command should have new-name flag")
	assert.True(t, updateAppCmd.Flags().Lookup(constants.OrgFlagName) != nil, "update app command should have organization flag")
	assert.True(t, updateAppCmd.Flags().Lookup(constants.OutputFlagName) != nil, "update app command should have output flag")
} 