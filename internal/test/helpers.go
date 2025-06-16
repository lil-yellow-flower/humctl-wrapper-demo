// Package test provides testing utilities for the humctl-wrapper-demo application.
// It includes helpers for executing commands and setting up mock clients.
//
// Testing Strategy:
// 1. Command Execution:
//    - Commands are tested in their actual hierarchy using the root command
//    - Output is captured and verified against expected results
//    - Error conditions are properly tested
//
// 2. Mock Client:
//    - A mock client is used to simulate Humanitec API responses
//    - The mock client can be configured to return specific responses or errors
//    - This allows testing different scenarios without making actual API calls
//
// 3. Test Structure:
//    - Each command has two types of tests:
//      a. Command Configuration: Tests the command's setup (flags, usage, etc.)
//      b. Command Execution: Tests the command's behavior with different inputs
//
// 4. Best Practices:
//    - Tests are independent and don't rely on external state
//    - Mock client is reset between tests
//    - Command output is properly captured and verified
//    - Error conditions are properly tested
package test

import (
	"bytes"
	"testing"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/spf13/cobra"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/config"
	"github.com/spf13/pflag"
)

// CommonTestArgs are common arguments for testing
var CommonTestArgs = []string{}

// ExecuteCommand executes a cobra command with the given arguments and flags
func ExecuteCommand(t *testing.T, root *cobra.Command, cmd *cobra.Command, args []string, flags map[string]string) (string, error) {
	t.Helper()

	// Create separate buffers for stdout and stderr
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	// Create a new root command for each test to ensure clean state
	testRoot := &cobra.Command{
		Use:   "test",
		Short: "Test root command",
	}
	testRoot.SetOut(stdout)
	testRoot.SetErr(stderr)

	// Create fresh copies of the commands
	freshRoot := &cobra.Command{
		Use:   root.Use,
		Short: root.Short,
	}
	freshCmd := &cobra.Command{
		Use:   cmd.Use,
		Short: cmd.Short,
		RunE:  cmd.RunE,
	}

	// Add the fresh root command to the test root
	testRoot.AddCommand(freshRoot)

	// Set default output format if not specified
	if _, ok := flags[constants.OutputFlagName]; !ok {
		flags[constants.OutputFlagName] = constants.DefaultOutputFormat
	}

	// Copy flags from original command to fresh command
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		// Create a new flag with the same properties
		newFlag := pflag.Flag{
			Name:      flag.Name,
			Shorthand: flag.Shorthand,
			Usage:     flag.Usage,
			Value:     flag.Value,
			DefValue:  flag.DefValue,
			Changed:   flag.Changed,
		}
		freshCmd.Flags().AddFlag(&newFlag)
	})

	// Copy required flags
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		// Check if the flag was required in the original command
		if cmd.Flags().Lookup(flag.Name) != nil && cmd.Flags().Lookup(flag.Name).Annotations != nil {
			if _, ok := cmd.Flags().Lookup(flag.Name).Annotations[cobra.BashCompOneRequiredFlag]; ok {
				freshCmd.MarkFlagRequired(flag.Name)
			}
		}
	})

	// Add the fresh command to the fresh root
	freshRoot.AddCommand(freshCmd)

	// Set flags on the fresh command
	for name, value := range flags {
		if err := freshCmd.Flags().Set(name, value); err != nil {
			t.Fatalf("Failed to set flag %s: %v", name, err)
		}
	}

	// Set arguments, prepending the command name
	cmdArgs := append([]string{freshRoot.Use, freshCmd.Use}, args...)
	testRoot.SetArgs(cmdArgs)

	// Set a dummy token to pass validation
	config.SetConfig(config.Config{
		HumanitecToken: "test-token",
		DefaultOutput:  flags[constants.OutputFlagName],
	})

	// Execute the command
	err := testRoot.Execute()

	// Return only the stdout content
	return stdout.String(), err
}

// SetupMockClient sets up a mock client for testing
func SetupMockClient(t *testing.T, mockClient *MockClient) {
	t.Helper()
	humanitec.SetClientFactory(&MockClientFactory{mockClient})
}

// MockClientFactory is a mock implementation of ClientFactory
type MockClientFactory struct {
	client *MockClient
}

// NewClient returns the mock client
func (f *MockClientFactory) NewClient(token, org string) humanitec.Client {
	return f.client
}

// MockClient is a mock implementation of Client
type MockClient struct {
	App   *humanitec.App
	Apps  []humanitec.App
	Error error
}

// GetApps returns the mock apps
func (c *MockClient) GetApps() ([]humanitec.App, error) {
	if c.Error != nil {
		return nil, c.Error
	}
	return c.Apps, nil
}

// GetApp returns the mock app
func (c *MockClient) GetApp(name string) (*humanitec.App, error) {
	if c.Error != nil {
		return nil, c.Error
	}
	return c.App, nil
}

// CreateApp returns the mock app
func (c *MockClient) CreateApp(id string, name string, skipEnvCreation bool) (*humanitec.App, error) {
	if c.Error != nil {
		return nil, c.Error
	}
	return c.App, nil
}

// DeleteApp returns the mock error
func (c *MockClient) DeleteApp(name string) error {
	return c.Error
}

// UpdateApp returns the mock app
func (c *MockClient) UpdateApp(oldName string, newName string) (*humanitec.App, error) {
	if c.Error != nil {
		return nil, c.Error
	}
	// Create a copy of the mock app with the new name
	updatedApp := *c.App
	updatedApp.Name = newName
	return &updatedApp, nil
} 