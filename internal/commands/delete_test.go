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

// mockDeleteClient implements the humanitec.Client interface for testing
type mockDeleteClient struct {
	error error
}

func (m *mockDeleteClient) GetApps() ([]humanitec.App, error) {
	return nil, nil
}

func (m *mockDeleteClient) CreateApp(name string, skipEnvCreation bool) (*humanitec.App, error) {
	return nil, nil
}

func (m *mockDeleteClient) DeleteApp(name string) error {
	return m.error
}

// setupTestDeleteCommand creates a test command with the given mock client
func setupTestDeleteCommand(t *testing.T, mockClient *mockDeleteClient) *cobra.Command {
	// Create root command
	rootCmd := &cobra.Command{
		Use:   "humctl",
		Short: "Test command",
	}

	// Create delete command
	deleteCmd := &cobra.Command{
		Use:   constants.DeleteCmdUse,
		Short: constants.DeleteCmdShort,
	}

	// Create app command
	appCmd := &cobra.Command{
		Use:   constants.DeleteAppCmdUse,
		Short: constants.DeleteAppCmdShort,
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

			// Delete application
			err = mockClient.DeleteApp(name)
			if err != nil {
				return fmt.Errorf(constants.ErrDeleteApp, err)
			}

			// Format and print output
			formatted, err := output.FormatMessage("Application successfully deleted", format)
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
	appCmd.Flags().StringP(constants.OrgFlagName, constants.OrgFlagShort, "", fmt.Sprintf("Humanitec organization ID (defaults to %s environment variable)", constants.HumanitecOrg))
	appCmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, constants.DefaultOutputFormat, constants.OutputFlagHelp)

	// Add commands to hierarchy
	deleteCmd.AddCommand(appCmd)
	rootCmd.AddCommand(deleteCmd)

	return rootCmd
}

func TestDeleteApp(t *testing.T) {
	// Setup test config
	configFile := testutil.SetupTestConfig(t)
	defer os.Remove(configFile)

	// Create a mock client
	mockClient := &mockDeleteClient{}

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
			args:           []string{"delete", "app", "--name", "test-app"},
			expectedOutput: "Application successfully deleted\n",
		},
		{
			name:           "json format",
			args:           []string{"delete", "app", "--name", "test-app", "--output", "json"},
			expectedOutput: "{\n  \"message\": \"Application successfully deleted\"\n}\n",
		},
		{
			name:           "yaml format",
			args:           []string{"delete", "app", "--name", "test-app", "--output", "yaml"},
			expectedOutput: "message: Application successfully deleted\n",
		},
		{
			name:           "successful delete with short flags",
			args:           []string{"delete", "app", "-n", "test-app", "-o", "json"},
			expectedOutput: "{\n  \"message\": \"Application successfully deleted\"\n}\n",
		},
		{
			name:           "successful delete with organization flag",
			args:           []string{"delete", "app", "--name", "test-app", "--org", "test-org"},
			expectedOutput: "Application successfully deleted\n",
		},
		{
			name:           "successful delete with organization short flag",
			args:           []string{"delete", "app", "-n", "test-app", "-g", "test-org"},
			expectedOutput: "Application successfully deleted\n",
		},
		{
			name:        "missing name flag",
			args:        []string{"delete", "app"},
			expectError: true,
		},
		{
			name:        "invalid output format",
			args:        []string{"delete", "app", "--name", "test-app", "--output", "invalid"},
			expectError: true,
		},
		{
			name:        "api error",
			args:        []string{"delete", "app", "--name", "test-app"},
			expectError: true,
			mockError:   fmt.Errorf("API error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set mock error if specified
			mockClient.error = tt.mockError

			// Create test command
			cmd := setupTestDeleteCommand(t, mockClient)

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

func TestDeleteCommandStructure(t *testing.T) {
	// Test delete command
	assert.Equal(t, constants.DeleteCmdUse, deleteCmd.Use, "delete command should have correct use")
	assert.Equal(t, constants.DeleteCmdShort, deleteCmd.Short, "delete command should have correct short description")

	// Test delete app command
	assert.Equal(t, constants.DeleteAppCmdUse, deleteAppCmd.Use, "delete app command should have correct use")
	assert.Equal(t, constants.DeleteAppCmdShort, deleteAppCmd.Short, "delete app command should have correct short description")

	// Test flags
	assert.True(t, deleteAppCmd.Flags().HasFlags(), "delete app command should have flags")
	assert.True(t, deleteAppCmd.Flags().Lookup(constants.NameFlagName) != nil, "delete app command should have name flag")
	assert.True(t, deleteAppCmd.Flags().Lookup(constants.OrgFlagName) != nil, "delete app command should have organization flag")
	assert.True(t, deleteAppCmd.Flags().Lookup(constants.OutputFlagName) != nil, "delete app command should have output flag")
} 