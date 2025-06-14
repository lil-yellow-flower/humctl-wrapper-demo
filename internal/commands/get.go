package commands

import (
	"fmt"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/output"
	"github.com/spf13/cobra"
)

var (
	// getCmd represents the get command
	getCmd = &cobra.Command{
		Use:   constants.GetCmdUse,
		Short: constants.GetCmdShort,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get format from flag or config
			formatStr, err := cmd.Flags().GetString(constants.OutputFlagName)
			if err != nil {
				return fmt.Errorf(constants.ErrInvalidOutputFormat, err)
			}

			// If format not specified, use default from config
			if formatStr == "" {
				formatStr = GetConfig().DefaultOutput
			}

			// Validate format
			format, err := output.ValidateFormat(formatStr)
			if err != nil {
				return fmt.Errorf(constants.ErrInvalidOutputFormat, err)
			}

			// Get applications
			client := humanitec.NewClient(GetConfig().HumanitecToken, GetConfig().HumanitecOrg)
			apps, err := client.GetApps()
			if err != nil {
				return fmt.Errorf(constants.ErrGetApps, err)
			}

			// Format and print output
			formatted, err := output.FormatApps(apps, format)
			if err != nil {
				return fmt.Errorf(constants.ErrFormatOutput, err)
			}
			fmt.Fprint(cmd.OutOrStdout(), formatted)

			return nil
		},
	}

	getAppCmd = &cobra.Command{
		Use:   constants.GetAppCmdUse,
		Short: constants.GetAppCmdShort,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get format from flag or config
			formatStr, err := cmd.Flags().GetString(constants.OutputFlagName)
			if err != nil {
				return fmt.Errorf(constants.ErrInvalidOutputFormat, err)
			}

			// If format not specified, use default from config
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

			// Get application
			app, err := client.GetApp(name)
			if err != nil {
				return fmt.Errorf(constants.ErrGetApp, err)
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
)

func init() {
	rootCmd.AddCommand(getCmd)

	// Add output format flag
	getCmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, "", constants.OutputFlagHelp)

	// Add flags to get app command
	getAppCmd.Flags().StringP(constants.NameFlagName, constants.NameFlagShort, "", constants.NameFlagHelp)
	getAppCmd.MarkFlagRequired(constants.NameFlagName)

	getAppCmd.Flags().StringP(constants.OrgFlagName, constants.OrgFlagShort, "", fmt.Sprintf("Humanitec organization ID (defaults to %s environment variable)", constants.HumanitecOrg))
	getAppCmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, "", constants.OutputFlagHelp)

	// Add commands to hierarchy
	getCmd.AddCommand(getAppCmd)
}