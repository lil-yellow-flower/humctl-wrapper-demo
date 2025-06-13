package commands

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/output"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Mock client for testing
type mockClient struct {
	apps  []humanitec.App
	error error
}

func (m *mockClient) ListApps() ([]humanitec.App, error) {
	return m.apps, m.error
}

// Mock client factory
func mockClientFactory(token, org string) humanitec.Client {
	return &mockClient{}
}

// setupTestCommand creates a new command hierarchy for testing
func setupTestCommand(t *testing.T, mockApps []humanitec.App, mockError error) (*cobra.Command, *bytes.Buffer) {
	// Create output buffer
	outBuf := new(bytes.Buffer)

	// Create root command
	rootCmd := &cobra.Command{
		Use:   constants.RootCmdUse,
		Short: "A command line interface wrapper for Humanitec platform",
		// Disable default error printing
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(outBuf)

	// Create get command
	getCmd := &cobra.Command{
		Use:   constants.GetCmdUse,
		Short: "Get resources from Humanitec platform",
		// Disable default error printing
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	getCmd.SetOut(outBuf)
	getCmd.SetErr(outBuf)

	// Create get apps command
	getAppsCmd := &cobra.Command{
		Use:   constants.GetAppsCmdUse,
		Short: "Get applications from Humanitec platform",
		// Disable default error printing
		SilenceErrors: true,
		SilenceUsage:  true,
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

			// Get org from flag or environment
			org, err := cmd.Flags().GetString(constants.OrgFlagName)
			if err != nil {
				return fmt.Errorf(constants.ErrInvalidOrgFlag, err)
			}

			// If org not provided, use environment variable
			if org == "" {
				org = os.Getenv(constants.EnvHumanitecOrg)
				if org == "" {
					return fmt.Errorf(constants.ErrMissingOrg)
				}
			}

			// Get token from environment
			token := os.Getenv(constants.EnvHumanitecToken)
			if token == "" {
				return fmt.Errorf(constants.ErrMissingToken)
			}

			// Use mock client
			c := &mockClient{
				apps:  mockApps,
				error: mockError,
			}

			// Get applications
			apps, err := c.ListApps()
			if err != nil {
				return fmt.Errorf(constants.ErrGetApps, err)
			}

			// Format and print output
			formatted, err := output.FormatApps(apps, format)
			if err != nil {
				return fmt.Errorf(constants.ErrFormatOutput, err)
			}
			fmt.Fprintln(outBuf, formatted)

			return nil
		},
	}
	getAppsCmd.SetOut(outBuf)
	getAppsCmd.SetErr(outBuf)

	// Add flags to get apps command
	getAppsCmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, constants.DefaultOutputFormat, "Output format (table|json|yaml)")
	getAppsCmd.Flags().String(constants.OrgFlagName, "", fmt.Sprintf("Humanitec organization ID (defaults to %s environment variable)", constants.EnvHumanitecOrg))

	// Add commands to hierarchy
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getAppsCmd)

	return rootCmd, outBuf
}

func TestGetApps(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		description    string
		args           []string
		envVars        map[string]string
		mockApps       []humanitec.App
		mockError      error
		expectedError  bool
		expectedOutput string
		expectedErrorMsg string
	}{
		{
			name:        "table_format",
			description: "should output applications in table format with default settings",
			args:        []string{constants.GetAppsCmdUse},
			envVars: map[string]string{
				constants.EnvHumanitecToken: "test-token",
				constants.EnvHumanitecOrg:   "test-org",
			},
			mockApps: []humanitec.App{
				{ID: "app1", Name: "Application 1"},
				{ID: "app2", Name: "Application 2"},
			},
			expectedError:  false,
			expectedOutput: "NAME\tID\n----\t--\nApplication 1\tapp1\nApplication 2\tapp2\n\n",
		},
		{
			name:        "json_format",
			description: "should output applications in JSON format when specified",
			args:        []string{constants.GetAppsCmdUse, "--output", "json"},
			envVars: map[string]string{
				constants.EnvHumanitecToken: "test-token",
				constants.EnvHumanitecOrg:   "test-org",
			},
			mockApps: []humanitec.App{
				{ID: "app1", Name: "Application 1"},
			},
			expectedError:  false,
			expectedOutput: "[\n  {\n    \"id\": \"app1\",\n    \"name\": \"Application 1\"\n  }\n]\n",
		},
		{
			name:        "yaml_format",
			description: "should output applications in YAML format when specified",
			args:        []string{constants.GetAppsCmdUse, "--output", "yaml"},
			envVars: map[string]string{
				constants.EnvHumanitecToken: "test-token",
				constants.EnvHumanitecOrg:   "test-org",
			},
			mockApps: []humanitec.App{
				{ID: "app1", Name: "Application 1"},
			},
			expectedError:  false,
			expectedOutput: "- id: app1\n  name: Application 1\n\n",
		},
		{
			name:        "invalid_format",
			description: "should return error for invalid output format",
			args:        []string{constants.GetAppsCmdUse, "--output", "invalid"},
			envVars: map[string]string{
				constants.EnvHumanitecToken: "test-token",
				constants.EnvHumanitecOrg:   "test-org",
			},
			expectedError: true,
			expectedErrorMsg: "invalid output format: unsupported output format: invalid. Supported formats: table, json, yaml",
		},
		{
			name:        "api_error",
			description: "should handle API errors gracefully",
			args:        []string{constants.GetAppsCmdUse},
			envVars: map[string]string{
				constants.EnvHumanitecToken: "test-token",
				constants.EnvHumanitecOrg:   "test-org",
			},
			mockError:     assert.AnError,
			expectedError: true,
			expectedErrorMsg: "failed to get applications: assert.AnError general error for testing",
		},
		{
			name:        "missing_token",
			description: "should require HUMANITEC_TOKEN environment variable",
			args:        []string{constants.GetAppsCmdUse},
			envVars: map[string]string{
				constants.EnvHumanitecOrg: "test-org",
			},
			expectedError: true,
			expectedErrorMsg: "HUMANITEC_TOKEN environment variable is required",
		},
		{
			name:        "missing_org",
			description: "should require HUMANITEC_ORG environment variable",
			args:        []string{constants.GetAppsCmdUse},
			envVars: map[string]string{
				constants.EnvHumanitecToken: "test-token",
			},
			expectedError: true,
			expectedErrorMsg: "HUMANITEC_ORG environment variable is required",
		},
		{
			name:        "org_flag",
			description: "should use org from flag when provided",
			args:        []string{constants.GetAppsCmdUse, "--org", "different-org"},
			envVars: map[string]string{
				constants.EnvHumanitecToken: "test-token",
				constants.EnvHumanitecOrg:   "test-org",
			},
			mockApps: []humanitec.App{
				{ID: "app1", Name: "Application 1"},
			},
			expectedError:  false,
			expectedOutput: "NAME\tID\n----\t--\nApplication 1\tapp1\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			// Set up command
			rootCmd, outBuf := setupTestCommand(t, tt.mockApps, tt.mockError)

			// Set command arguments
			rootCmd.SetArgs(append([]string{constants.GetCmdUse}, tt.args...))

			// Execute command
			err := rootCmd.Execute()

			// Check error
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
					t.Logf("Test case: %s", tt.description)
					t.Logf("Command args: %v", tt.args)
					t.Logf("Environment variables: %v", tt.envVars)
					t.Logf("Command output: %s", outBuf.String())
					return
				}
				if err.Error() != tt.expectedErrorMsg {
					t.Errorf("error message = %v, want %v", err.Error(), tt.expectedErrorMsg)
					t.Logf("Test case: %s", tt.description)
					t.Logf("Command args: %v", tt.args)
					t.Logf("Environment variables: %v", tt.envVars)
					t.Logf("Command output: %s", outBuf.String())
					return
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					t.Logf("Test case: %s", tt.description)
					t.Logf("Command args: %v", tt.args)
					t.Logf("Environment variables: %v", tt.envVars)
					t.Logf("Command output: %s", outBuf.String())
					return
				}
				if got := outBuf.String(); got != tt.expectedOutput {
					t.Errorf("output = %v, want %v", got, tt.expectedOutput)
					t.Logf("Test case: %s", tt.description)
					t.Logf("Command args: %v", tt.args)
					t.Logf("Environment variables: %v", tt.envVars)
					return
				}
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