package commands

import (
	"fmt"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/output"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   constants.UpdateCmdUse,
	Short: constants.UpdateCmdShort,
}

var updateAppCmd = &cobra.Command{
	Use:   constants.UpdateAppCmdUse,
	Short: constants.UpdateAppCmdShort,
	Long: `Update an existing application in the organization.
Currently supports updating the application name while preserving all other settings and configurations.`,
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
			return fmt.Errorf(constants.ErrInvalidOrgFlag, err)
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

		// Get new name from flag
		newName, err := cmd.Flags().GetString(constants.NewNameFlagName)
		if err != nil {
			return fmt.Errorf(constants.ErrInvalidNewName, err)
		}

		// Update application
		_, err = client.UpdateApp(name, newName)
		if err != nil {
			return fmt.Errorf(constants.ErrUpdateApp, err)
		}

		// Format and print output
		formatted, err := output.FormatMessage(constants.SuccessAppUpdated, format)
		if err != nil {
			return fmt.Errorf(constants.ErrFormatOutput, err)
		}
		fmt.Fprint(cmd.OutOrStdout(), formatted)

		return nil
	},
}

func init() {
	// Add flags to update app command
	updateAppCmd.Flags().StringP(constants.NameFlagName, constants.NameFlagShort, "", constants.NameFlagHelp)
	updateAppCmd.MarkFlagRequired(constants.NameFlagName)

	updateAppCmd.Flags().StringP(constants.NewNameFlagName, constants.NewNameFlagShort, "", constants.NewNameFlagHelp)
	updateAppCmd.MarkFlagRequired(constants.NewNameFlagName)

	updateAppCmd.Flags().StringP(constants.OrgFlagName, constants.OrgFlagShort, "", fmt.Sprintf("Humanitec organization ID (defaults to %s environment variable)", constants.HumanitecOrg))
	updateAppCmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, "", constants.OutputFlagHelp)

	// Add commands to hierarchy
	updateCmd.AddCommand(updateAppCmd)
	rootCmd.AddCommand(updateCmd)
} 