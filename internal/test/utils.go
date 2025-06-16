package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
)

// SetupTestConfig creates a temporary config.yaml file for testing
func SetupTestConfig(t *testing.T) string {
	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Write test config
	config := fmt.Sprintf(`%s: "test-token"
%s: "test-org"
default_output: "%s"
logging:
  level: "info"
  format: "text"
  output: "stdout"
  file: "logs/humctl-wrapper.log"`, constants.HumanitecToken, constants.HumanitecOrg, constants.DefaultOutputFormat)

	if _, err := tmpFile.WriteString(config); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close config file: %v", err)
	}

	return tmpFile.Name()
} 