package commands

import (
	"fmt"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/output"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   constants.AddCmdUse,
	Short: constants.AddCmdShort,
}

var addAppCmd = &cobra.Command{
	Use:   constants.AddAppCmdUse,
	Short: constants.AddAppCmdShort,
	Long: `Add a new application to the organization with the specified name.
The application ID will be automatically generated from the name by converting it to lowercase and replacing spaces with hyphens.`,
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
		org, err := cmd.Flags().GetString(constants.OrgFlagName)
		if err != nil {
			return fmt.Errorf("failed to get organization ID: %w", err)
		}
		if org == "" {
			org = GetConfig().HumanitecOrg
		}

		// Get token from config
		token := GetConfig().HumanitecToken

		// Create Humanitec client
		client := humanitec.NewClient(token, org)

		// Get name from flag
		name, err := cmd.Flags().GetString(constants.NameFlagName)
		if err != nil {
			return fmt.Errorf(constants.ErrInvalidName, err)
		}

		// Get skip environment creation flag
		skipEnvCreation, err := cmd.Flags().GetBool(constants.SkipEnvCreationFlagName)
		if err != nil {
			return fmt.Errorf(constants.ErrInvalidSkipEnvCreation, err)
		}

		// Add application
		app, err := client.CreateApp(name, skipEnvCreation)
		if err != nil {
			return fmt.Errorf(constants.ErrAddApp, err)
		}

		// Format and print output
		formatted, err := output.FormatApp(app, format)
		if err != nil {
			return fmt.Errorf(constants.ErrFormatOutput, err)
		}
		fmt.Fprint(cmd.OutOrStdout(), formatted)

		return nil
	},
}

func init() {
	// Add flags to add app command
	addAppCmd.Flags().StringP(constants.NameFlagName, constants.NameFlagShort, "", constants.NameFlagHelp)
	addAppCmd.MarkFlagRequired(constants.NameFlagName)

	addAppCmd.Flags().BoolP(constants.SkipEnvCreationFlagName, constants.SkipEnvCreationFlagShort, false, constants.SkipEnvCreationFlagHelp)
	addAppCmd.Flags().StringP(constants.OrgFlagName, constants.OrgFlagShort, "", fmt.Sprintf("Humanitec organization ID (defaults to %s environment variable)", constants.HumanitecOrg))
	addAppCmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, "", constants.OutputFlagHelp)

	// Add commands to hierarchy
	addCmd.AddCommand(addAppCmd)
	rootCmd.AddCommand(addCmd)
} 