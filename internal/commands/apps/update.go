package apps

import (
	"fmt"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/config"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/output"
	"github.com/spf13/cobra"
)

var (
	// Subcommand for updating apps
	update = &cobra.Command{
		Use:   constants.AppCmdUse,
		Short: constants.AppCmdShort,
		Long: `Update an existing application in the organization.
Currently supports updating the application name while preserving all other settings and configurations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flags
			id, err := cmd.Flags().GetString(constants.IDFlagName)
			if err != nil {
				return fmt.Errorf("failed to get id flag: %w", err)
			}

			name, err := cmd.Flags().GetString(constants.NameFlagName)
			if err != nil {
				return fmt.Errorf("failed to get name flag: %w", err)
			}

			outputFormatStr, err := cmd.Flags().GetString(constants.OutputFlagName)
			if err != nil {
				return fmt.Errorf("failed to get output format flag: %w", err)
			}

			// Validate output format
			outputFormat, err := output.ValidateFormat(outputFormatStr)
			if err != nil {
				return fmt.Errorf("invalid output format: %w", err)
			}

			// Get organization ID from config
			org := config.GetConfig().HumanitecOrg
			token := config.GetConfig().HumanitecToken

			// Create Humanitec client
			client := humanitec.NewClient(token, org)

			// Update app
			app, err := client.UpdateApp(id, name)
			if err != nil {
				return fmt.Errorf("failed to update app: %w", err)
			}

			// Print output
			formatted, err := output.FormatApp(app, outputFormat)
			if err != nil {
				return fmt.Errorf("failed to format output: %w", err)
			}
			fmt.Fprint(cmd.OutOrStdout(), formatted)

			return nil
		},
	}
)

func init() {
	// Add common flags
	CommonFlagSet()(update)

	// Add command-specific flags
	update.Flags().StringP(constants.IDFlagName, constants.IDFlagShort, "", constants.IDFlagHelp)
	update.Flags().StringP(constants.NameFlagName, constants.NameFlagShort, "", constants.NameFlagHelp)

	// Mark required flags
	update.MarkFlagRequired(constants.IDFlagName)
	update.MarkFlagRequired(constants.NameFlagName)
} 