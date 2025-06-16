package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"gopkg.in/yaml.v3"
)

// Format represents the output format
type Format string

const (
	// FormatTable represents table output format
	FormatTable Format = "table"
	// FormatJSON represents JSON output format
	FormatJSON Format = "json"
	// FormatYAML represents YAML output format
	FormatYAML Format = "yaml"
)

// ValidateFormat validates if the given format is supported
func ValidateFormat(format string) (Format, error) {
	switch Format(strings.ToLower(format)) {
	case FormatTable, FormatJSON, FormatYAML:
		return Format(strings.ToLower(format)), nil
	default:
		return FormatTable, fmt.Errorf("unsupported output format: %s. Supported formats: table, json, yaml", format)
	}
}

// FormatApps formats a list of applications in the specified format.
// JSON and Table formats include a trailing newline, while YAML format
// uses the newline provided by the YAML marshaler.
func FormatApps(apps []humanitec.App, format Format) (string, error) {
	switch format {
	case FormatJSON:
		data, err := json.MarshalIndent(apps, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal to JSON: %w", err)
		}
		return string(data) + "\n", nil

	case FormatYAML:
		data, err := yaml.Marshal(apps)
		if err != nil {
			return "", fmt.Errorf("failed to marshal to YAML: %w", err)
		}
		return string(data), nil

	case FormatTable:
		var sb strings.Builder
		sb.WriteString("NAME\tID\n")
		sb.WriteString("----\t--\n")
		for _, app := range apps {
			sb.WriteString(fmt.Sprintf("%s\t%s\n", app.Name, app.ID))
		}
		return sb.String(), nil

	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// FormatApp formats a single application in the specified format.
// JSON and Table formats include a trailing newline, while YAML format
// uses the newline provided by the YAML marshaler.
func FormatApp(app *humanitec.App, format Format) (string, error) {
	switch format {
	case FormatJSON:
		data, err := json.MarshalIndent(app, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal to JSON: %w", err)
		}
		return string(data) + "\n", nil

	case FormatYAML:
		data, err := yaml.Marshal(app)
		if err != nil {
			return "", fmt.Errorf("failed to marshal to YAML: %w", err)
		}
		return string(data), nil

	case FormatTable:
		var sb strings.Builder
		sb.WriteString("NAME\tID\n")
		sb.WriteString("----\t--\n")
		sb.WriteString(fmt.Sprintf("%s\t%s\n", app.Name, app.ID))
		return sb.String(), nil

	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// FormatMessage formats a simple message in the specified format.
// JSON and Table formats include a trailing newline, while YAML format
// uses the newline provided by the YAML marshaler.
func FormatMessage(message string, format Format) (string, error) {
	switch format {
	case FormatJSON:
		data, err := json.MarshalIndent(map[string]string{"message": message}, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal to JSON: %w", err)
		}
		return string(data) + "\n", nil

	case FormatYAML:
		data, err := yaml.Marshal(map[string]string{"message": message})
		if err != nil {
			return "", fmt.Errorf("failed to marshal to YAML: %w", err)
		}
		// YAML marshaling adds a newline, so we don't need to add another one
		return string(data), nil

	case FormatTable:
		return message + "\n", nil

	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
} 