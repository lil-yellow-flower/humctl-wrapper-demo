package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/output"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var getCmd = &cobra.Command{
	Use:   constants.GetCmdUse,
	Short: "Get resources from Humanitec platform",
}

var getAppsCmd = &cobra.Command{
	Use:   constants.GetAppsCmdUse,
	Short: "Get applications from Humanitec platform",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get format from flag or config
		formatStr, err := cmd.Flags().GetString(constants.OutputFlagName)
		if err != nil {
			return fmt.Errorf(constants.ErrInvalidOutputFormat, err)
		}

		// If no format specified, use default from config
		if formatStr == "" {
			formatStr = GetConfig().DefaultOutput
		}

		// Validate format
		format, err := output.ValidateFormat(formatStr)
		if err != nil {
			return fmt.Errorf(constants.ErrInvalidOutputFormat, err)
		}

		// Get organization ID from flag or config
		org := cmd.Flag("org").Value.String()
		if org == "" {
			org = GetConfig().HumanitecOrg
		}

		// Get token from config
		token := GetConfig().HumanitecToken

		// Create Humanitec client
		client := humanitec.NewClient(token, org)

		// Get applications
		apps, err := client.GetApps()
		if err != nil {
			return fmt.Errorf(constants.ErrGetApps, err)
		}

		// Format and print output
		formatted, err := output.FormatApps(apps, format)
		if err != nil {
			return fmt.Errorf(constants.ErrFormatOutput, err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), formatted)

		return nil
	},
}

// printJSON prints the applications in JSON format
func printJSON(apps []humanitec.App) error {
	data, err := json.MarshalIndent(apps, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// printYAML prints the applications in YAML format
func printYAML(apps []humanitec.App) error {
	data, err := yaml.Marshal(apps)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// printTable prints the applications in table format
func printTable(apps []humanitec.App) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tName")
	for _, app := range apps {
		fmt.Fprintf(w, "%s\t%s\n", app.ID, app.Name)
	}
	return w.Flush()
}

func init() {
	// Add flags to get apps command
	getAppsCmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, "", "Output format (table|json|yaml)")
	getAppsCmd.Flags().StringP(constants.OrgFlagName, constants.OrgFlagShort, "", fmt.Sprintf("Humanitec organization ID (defaults to %s environment variable)", constants.HumanitecOrg))

	// Add commands to hierarchy
	getCmd.AddCommand(getAppsCmd)
	rootCmd.AddCommand(getCmd)
}