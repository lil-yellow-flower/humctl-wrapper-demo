package commands

import (
	"fmt"

	"github.com/mathi-ma51zaw/humctl-wrapper-demo/internal/humanitec"
	"github.com/mathi-ma51zaw/humctl-wrapper-demo/internal/output"
	"github.com/spf13/cobra"
)

var (
	// Parent command
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get Humanitec resources",
		Long: `Get Humanitec resources such as applications, environments, and deployments.
For example:
  humctl-wrapper get apps
  humctl-wrapper get apps --output json
  humctl-wrapper get apps --output yaml`,
	}

	// Subcommand for listing all applications
	getAppsCmd = &cobra.Command{
		Use:     "apps",
		Aliases: []string{"app"},
		Short:   "Get applications",
		Long: `Get applications in your Humanitec organization.
Examples:
  # List all applications
  humctl-wrapper get apps

  # Get a specific application
  humctl-wrapper get apps my-app

  # Get applications in JSON format
  humctl-wrapper get apps --output json

  # Get applications in YAML format
  humctl-wrapper get apps --output yaml`,
		RunE: runGetApps,
	}
)

func init() {
	// Register subcommands under the parent command
	getCmd.AddCommand(getAppsCmd)
	rootCmd.AddCommand(getCmd)

	// Add flags
	getCmd.PersistentFlags().StringP("output", "o", "table", "Output format (table, json, yaml)")
}

func runGetApps(cmd *cobra.Command, args []string) error {
	client := humanitec.NewClient()
	apps, err := client.ListApps()
	if err != nil {
		return fmt.Errorf("failed to list applications: %w", err)
	}

	// Get output format
	outputFormat, _ := cmd.Flags().GetString("output")
	format, err := output.ValidateFormat(outputFormat)
	if err != nil {
		return err
	}

	// Format and print output
	formatted, err := output.FormatApps(apps, format)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	fmt.Print(formatted)
	return nil
} 