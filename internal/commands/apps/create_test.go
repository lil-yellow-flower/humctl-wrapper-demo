package apps

import (
	"fmt"
	"testing"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/test"
	"github.com/stretchr/testify/assert"
)

var createTestCases = []struct {
	name           string
	args           []string
	flags          map[string]string
	expectedOutput string
	expectError    bool
	mockError      error
}{
	{
		name:           "create app with valid id and name - table format",
		args:           []string{"create"},
		flags:          map[string]string{constants.IDFlagName: "test-app", constants.NameFlagName: "Test App", "output": "table"},
		expectedOutput: "NAME\tID\n----\t--\nTest App\ttest-app\n",
		expectError:    false,
	},
	{
		name:           "create app with valid id and name - json format",
		args:           []string{"create"},
		flags:          map[string]string{constants.IDFlagName: "test-app", constants.NameFlagName: "Test App", "output": "json"},
		expectedOutput: "{\n  \"id\": \"test-app\",\n  \"name\": \"Test App\"\n}\n",
		expectError:    false,
	},
	{
		name:           "create app with valid id and name - yaml format",
		args:           []string{"create"},
		flags:          map[string]string{constants.IDFlagName: "test-app", constants.NameFlagName: "Test App", "output": "yaml"},
		expectedOutput: "id: test-app\nname: Test App\n",
		expectError:    false,
	},
	{
		name:           "create app with skip environment creation",
		args:           []string{"create"},
		flags:          map[string]string{constants.IDFlagName: "test-app", constants.NameFlagName: "Test App", "skip-env-creation": "true", "output": "table"},
		expectedOutput: "NAME\tID\n----\t--\nTest App\ttest-app\n",
		expectError:    false,
	},
	{
		name:           "missing id flag",
		args:           []string{"create"},
		flags:          map[string]string{constants.NameFlagName: "Test App", "output": "table"},
		expectedOutput: "required flag(s) \"id\" not set",
		expectError:    true,
	},
	{
		name:           "missing name flag",
		args:           []string{"create"},
		flags:          map[string]string{constants.IDFlagName: "test-app", "output": "table"},
		expectedOutput: "required flag(s) \"name\" not set",
		expectError:    true,
	},
	{
		name:           "invalid output format",
		args:           []string{"create"},
		flags:          map[string]string{constants.IDFlagName: "test-app", constants.NameFlagName: "Test App", "output": "invalid"},
		expectedOutput: "invalid output format: unsupported output format: invalid. Supported formats: table, json, yaml",
		expectError:    true,
	},
	{
		name:           "api error - duplicate id",
		args:           []string{"create"},
		flags:          map[string]string{constants.IDFlagName: "existing-app", constants.NameFlagName: "Test App", "output": "table"},
		expectedOutput: "failed to create app: application with id 'existing-app' already exists",
		expectError:    true,
		mockError:      fmt.Errorf("application with id 'existing-app' already exists"),
	},
	{
		name:           "api error - server error",
		args:           []string{"create"},
		flags:          map[string]string{constants.IDFlagName: "test-app", constants.NameFlagName: "Test App", "output": "table"},
		expectedOutput: "failed to create app: API error",
		expectError:    true,
		mockError:      fmt.Errorf("API error"),
	},
}

// TestCreateAppCommandExecution verifies the create app command's runtime behavior by testing
// various input combinations, output formats, and error conditions using a mock client.
func TestCreateAppCommandExecution(t *testing.T) {
	// Create a mock client
	mockClient := &test.MockClient{
		App: &humanitec.App{
			ID:   "test-app",
			Name: "Test App",
		},
	}

	// Set up the mock client
	test.SetupMockClient(t, mockClient)

	for _, tt := range createTestCases {
		t.Run(tt.name, func(t *testing.T) {
			// Set mock error if specified
			mockClient.Error = tt.mockError

			got, err := test.ExecuteCommand(t, create, create, tt.args, tt.flags)
			if (err != nil) != tt.expectError {
				t.Errorf("create.Execute() error = %v, wantErr %v", err, tt.expectError)
				return
			}
			if !tt.expectError && got != tt.expectedOutput {
				t.Errorf("create.Execute() = %v, want %v", got, tt.expectedOutput)
			}
			if tt.expectError && err != nil && err.Error() != tt.expectedOutput {
				t.Errorf("create.Execute() error = %v, want %v", err.Error(), tt.expectedOutput)
			}
		})
	}
}

// TestCreateAppCommandConfiguration verifies that the create app command is properly configured
// with the correct name, description, and required flags.
func TestCreateAppCommandConfiguration(t *testing.T) {
	// Test create command
	assert.Equal(t, constants.AppCmdUse, create.Use, "create command should have correct use")
	assert.Equal(t, constants.AppCmdShort, create.Short, "create command should have correct short description")

	// Test flags
	assert.True(t, create.Flags().HasFlags(), "create command should have flags")
	assert.True(t, create.Flags().Lookup(constants.IDFlagName) != nil, "create command should have id flag")
	assert.True(t, create.Flags().Lookup(constants.NameFlagName) != nil, "create command should have name flag")
	assert.True(t, create.Flags().Lookup(constants.OutputFlagName) != nil, "create command should have output flag")
	assert.True(t, create.Flags().Lookup(constants.SkipEnvCreationFlagName) != nil, "create command should have skip-env-creation flag")
} 