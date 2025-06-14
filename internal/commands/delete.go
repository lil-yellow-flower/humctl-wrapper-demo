package commands

import (
	"fmt"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/output"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   constants.DeleteCmdUse,
	Short: constants.DeleteCmdShort,
}

var deleteAppCmd = &cobra.Command{
	Use:   constants.DeleteAppCmdUse,
	Short: constants.DeleteAppCmdShort,
	Long: `Delete an Application and everything associated with it. This includes:
- Environments
- Deployment history on those Environments
- Any shared values and secrets associated

Deletions are currently irreversible.`,
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

		// Delete application
		err = client.DeleteApp(name)
		if err != nil {
			return fmt.Errorf(constants.ErrDeleteApp, err)
		}

		// Format and print output
		formatted, err := output.FormatMessage("Application successfully deleted", format)
		if err != nil {
			return fmt.Errorf(constants.ErrFormatOutput, err)
		}
		fmt.Fprint(cmd.OutOrStdout(), formatted)

		return nil
	},
}

func init() {
	// Add flags to delete app command
	deleteAppCmd.Flags().StringP(constants.NameFlagName, constants.NameFlagShort, "", constants.NameFlagHelp)
	deleteAppCmd.MarkFlagRequired(constants.NameFlagName)

	deleteAppCmd.Flags().StringP(constants.OrgFlagName, constants.OrgFlagShort, "", fmt.Sprintf("Humanitec organization ID (defaults to %s environment variable)", constants.HumanitecOrg))
	deleteAppCmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, "", constants.OutputFlagHelp)

	// Add commands to hierarchy
	deleteCmd.AddCommand(deleteAppCmd)
	rootCmd.AddCommand(deleteCmd)
} 