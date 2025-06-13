package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mathi-ma51zaw/humctl-wrapper-demo/internal/humanitec"
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

// FormatApps formats the applications in the specified format
func FormatApps(apps []humanitec.App, format Format) (string, error) {
	switch format {
	case FormatJSON:
		data, err := json.MarshalIndent(apps, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal to JSON: %w", err)
		}
		return string(data), nil

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