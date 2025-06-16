package apps

import (
	"fmt"
	"testing"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/test"
	"github.com/stretchr/testify/assert"
)

var updateTestCases = []struct {
	name           string
	args           []string
	flags          map[string]string
	expectedOutput string
	expectError    bool
	mockError      error
}{
	{
		name:           "table format",
		args:           []string{"update"},
		flags:          map[string]string{"id": "test-app", "name": "New App Name", "output": "table"},
		expectedOutput: "NAME\tID\n----\t--\nNew App Name\ttest-app\n",
		expectError:    false,
	},
	{
		name:           "json format",
		args:           []string{"update"},
		flags:          map[string]string{"id": "test-app", "name": "New App Name", "output": "json"},
		expectedOutput: "{\n  \"id\": \"test-app\",\n  \"name\": \"New App Name\"\n}\n",
		expectError:    false,
	},
	{
		name:           "yaml format",
		args:           []string{"update"},
		flags:          map[string]string{"id": "test-app", "name": "New App Name", "output": "yaml"},
		expectedOutput: "id: test-app\nname: New App Name\n",
		expectError:    false,
	},
	{
		name:           "missing id flag",
		args:           []string{"update"},
		flags:          map[string]string{"name": "New App Name"},
		expectedOutput: "",
		expectError:    true,
	},
	{
		name:           "missing name flag",
		args:           []string{"update"},
		flags:          map[string]string{"id": "test-app"},
		expectedOutput: "",
		expectError:    true,
	},
	{
		name:           "invalid output format",
		args:           []string{"update"},
		flags:          map[string]string{"id": "test-app", "name": "New App Name", "output": "invalid"},
		expectedOutput: "",
		expectError:    true,
	},
	{
		name:           "api error",
		args:           []string{"update"},
		flags:          map[string]string{"id": "error-app", "name": "New App Name", "output": "table"},
		expectedOutput: "",
		expectError:    true,
		mockError:      fmt.Errorf("API error"),
	},
}

// TestUpdateAppCommandExecution verifies the update app command's runtime behavior by testing
// various input combinations, output formats, and error conditions using a mock client.
func TestUpdateAppCommandExecution(t *testing.T) {
	// Create a mock client
	mockClient := &test.MockClient{
		App: &humanitec.App{
			ID:   "test-app",
			Name: "test-app",
		},
	}

	// Set up the mock client
	test.SetupMockClient(t, mockClient)

	for _, tt := range updateTestCases {
		t.Run(tt.name, func(t *testing.T) {
			// Set mock error if specified
			mockClient.Error = tt.mockError

			got, err := test.ExecuteCommand(t, update, update, tt.args, tt.flags)
			if (err != nil) != tt.expectError {
				t.Errorf("update.Execute() error = %v, wantErr %v", err, tt.expectError)
				return
			}
			if !tt.expectError && got != tt.expectedOutput {
				t.Errorf("update.Execute() = %v, want %v", got, tt.expectedOutput)
			}
		})
	}
}

// TestUpdateAppCommandConfiguration verifies that the update app command is properly configured
// with the correct name, description, and required flags.
func TestUpdateAppCommandConfiguration(t *testing.T) {
	// Test update command
	assert.Equal(t, constants.AppCmdUse, update.Use, "update command should have correct use")
	assert.Equal(t, constants.AppCmdShort, update.Short, "update command should have correct short description")

	// Test flags
	assert.True(t, update.Flags().HasFlags(), "update command should have flags")
	assert.True(t, update.Flags().Lookup(constants.IDFlagName) != nil, "update command should have id flag")
	assert.True(t, update.Flags().Lookup(constants.NameFlagName) != nil, "update command should have name flag")
	assert.True(t, update.Flags().Lookup(constants.OutputFlagName) != nil, "update command should have output flag")
} 