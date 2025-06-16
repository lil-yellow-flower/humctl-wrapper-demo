package apps

import (
	"fmt"
	"testing"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/test"
	"github.com/stretchr/testify/assert"
)

var deleteTestCases = []struct {
	name           string
	args           []string
	flags          map[string]string
	expectedOutput string
	expectError    bool
	mockError      error
}{
	{
		name:  "table format",
		args:  []string{"delete"},
		flags: map[string]string{
			constants.IDFlagName:    "test-app",
			constants.OrgFlagName:   "test-org",
			constants.OutputFlagName: constants.DefaultOutputFormat,
		},
		expectedOutput: "Application successfully deleted\n",
		expectError:    false,
	},
	{
		name:  "json format",
		args:  []string{"delete"},
		flags: map[string]string{
			constants.IDFlagName:    "test-app",
			constants.OrgFlagName:   "test-org",
			constants.OutputFlagName: "json",
		},
		expectedOutput: "{\n  \"message\": \"Application successfully deleted\"\n}\n",
		expectError:    false,
	},
	{
		name:  "yaml format",
		args:  []string{"delete"},
		flags: map[string]string{
			constants.IDFlagName:    "test-app",
			constants.OrgFlagName:   "test-org",
			constants.OutputFlagName: "yaml",
		},
		expectedOutput: "message: Application successfully deleted\n",
		expectError:    false,
	},
	{
		name:        "missing id flag",
		args:        []string{"delete"},
		flags:       map[string]string{},
		expectError: true,
	},
	{
		name:  "invalid output format",
		args:  []string{"delete"},
		flags: map[string]string{
			constants.IDFlagName:    "test-app",
			constants.OrgFlagName:   "test-org",
			constants.OutputFlagName: "invalid",
		},
		expectError: true,
	},
	{
		name:  "api error",
		args:  []string{"delete"},
		flags: map[string]string{
			constants.IDFlagName:  "test-app",
			constants.OrgFlagName: "test-org",
		},
		expectError: true,
		mockError:   fmt.Errorf("API error"),
	},
}

// TestDeleteAppCommandExecution verifies the delete app command's runtime behavior by testing
// various input combinations, output formats, and error conditions using a mock client.
func TestDeleteAppCommandExecution(t *testing.T) {
	// Create a mock client
	mockClient := &test.MockClient{
		App: &humanitec.App{
			ID:   "test-app",
			Name: "Test App",
		},
	}

	// Set up the mock client
	test.SetupMockClient(t, mockClient)

	for _, tt := range deleteTestCases {
		t.Run(tt.name, func(t *testing.T) {
			// Set mock error if specified
			mockClient.Error = tt.mockError

			got, err := test.ExecuteCommand(t, delete, delete, tt.args, tt.flags)
			if (err != nil) != tt.expectError {
				t.Errorf("delete.Execute() error = %v, wantErr %v", err, tt.expectError)
				return
			}
			if !tt.expectError && got != tt.expectedOutput {
				t.Errorf("delete.Execute() = %v, want %v", got, tt.expectedOutput)
			}
		})
	}
}

// TestDeleteAppCommandConfiguration verifies that the delete app command is properly configured
// with the correct name, description, and required flags.
func TestDeleteAppCommandConfiguration(t *testing.T) {
	// Test delete command
	assert.Equal(t, constants.AppCmdUse, delete.Use, "delete command should have correct use")
	assert.Equal(t, constants.AppCmdShort, delete.Short, "delete command should have correct short description")

	// Test flags
	assert.True(t, delete.Flags().HasFlags(), "delete command should have flags")
	assert.True(t, delete.Flags().Lookup(constants.IDFlagName) != nil, "delete command should have id flag")
	assert.True(t, delete.Flags().Lookup(constants.OutputFlagName) != nil, "delete command should have output flag")
} 