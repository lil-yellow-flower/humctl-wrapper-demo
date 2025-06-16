package apps

import (
	"testing"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestGetAppCommandExecution(t *testing.T) {
	// Test cases for get command
	getTestCases := []struct {
		name           string
		args           []string
		flags          map[string]string
		expectedOutput string
		expectedError  bool
		mockError      error
	}{
		{
			name:           "get single app - table format",
			args:           []string{},
			flags:          map[string]string{constants.IDFlagName: "test-app", constants.OutputFlagName: "table"},
			expectedOutput: "NAME\tID\n----\t--\ntest-app\ttest-app\n",
			expectedError:  false,
		},
		{
			name:           "get single app - json format",
			args:           []string{},
			flags:          map[string]string{constants.IDFlagName: "test-app", constants.OutputFlagName: "json"},
			expectedOutput: "{\n  \"id\": \"test-app\",\n  \"name\": \"test-app\"\n}\n",
			expectedError:  false,
		},
		{
			name:           "get single app - yaml format",
			args:           []string{},
			flags:          map[string]string{constants.IDFlagName: "test-app", constants.OutputFlagName: "yaml"},
			expectedOutput: "id: test-app\nname: test-app\n",
			expectedError:  false,
		},
		{
			name:           "list all apps - table format",
			args:           []string{},
			flags:          map[string]string{constants.OutputFlagName: "table"},
			expectedOutput: "NAME\tID\n----\t--\ntest-app\ttest-app\n",
			expectedError:  false,
		},
		{
			name:           "list all apps - json format",
			args:           []string{},
			flags:          map[string]string{constants.OutputFlagName: "json"},
			expectedOutput: "{\n  \"id\": \"test-app\",\n  \"name\": \"test-app\"\n}\n",
			expectedError:  false,
		},
		{
			name:           "list all apps - yaml format",
			args:           []string{},
			flags:          map[string]string{constants.OutputFlagName: "yaml"},
			expectedOutput: "id: test-app\nname: test-app\n",
			expectedError:  false,
		},
		{
			name:           "invalid output format",
			args:           []string{},
			flags:          map[string]string{constants.IDFlagName: "test-app", constants.OutputFlagName: "invalid"},
			expectedOutput: "Usage:\n  test apps apps [flags]\n\nFlags:\n  -h, --help            help for apps\n  -i, --id string       Application ID\n  -g, --org string      Humanitec organization ID (defaults to %s environment variable)\n  -o, --output string   Output format (table|json|yaml) (default \"table\")\n\n",
			expectedError:  true,
		},
		{
			name:           "api error",
			args:           []string{},
			flags:          map[string]string{constants.IDFlagName: "test-app", constants.OutputFlagName: "table"},
			expectedOutput: "Usage:\n  test apps apps [flags]\n\nFlags:\n  -h, --help            help for apps\n  -i, --id string       Application ID\n  -g, --org string      Humanitec organization ID (defaults to %s environment variable)\n  -o, --output string   Output format (table|json|yaml) (default \"table\")\n\n",
			expectedError:  true,
			mockError:      assert.AnError,
		},
	}

	for _, tc := range getTestCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock client
			mockClient := &test.MockClient{
				App: &humanitec.App{
					ID:   "test-app",
					Name: "test-app",
				},
				Apps: []humanitec.App{
					{
						ID:   "test-app",
						Name: "test-app",
					},
				},
				Error: tc.mockError,
			}

			// Set up the mock client
			test.SetupMockClient(t, mockClient)

			// Execute command
			output, err := test.ExecuteCommand(t, get, get, tc.args, tc.flags)

			// Check error
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Check output
			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}

func TestGetAppCommandConfiguration(t *testing.T) {
	// Test command configuration
	assert.Equal(t, constants.AppsCmdUse, get.Use)
	assert.Equal(t, constants.AppsCmdShort, get.Short)

	// Test required flags
	assert.True(t, get.Flags().Lookup(constants.OutputFlagName) != nil)
} 