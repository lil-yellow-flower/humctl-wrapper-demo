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
	// Subcommand for getting apps
	get = &cobra.Command{
		Use:   constants.AppsCmdUse,
		Short: constants.AppsCmdShort,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flags
			id, err := cmd.Flags().GetString(constants.IDFlagName)
			if err != nil {
				return fmt.Errorf("failed to get id flag: %w", err)
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

			// If ID is provided, get single app
			if id != "" {
				app, err := client.GetApp(id)
				if err != nil {
					return fmt.Errorf("failed to get app: %w", err)
				}

				// Print output
				formatted, err := output.FormatApp(app, outputFormat)
				if err != nil {
					return fmt.Errorf("failed to format output: %w", err)
				}
				fmt.Fprint(cmd.OutOrStdout(), formatted)

				return nil
			}

			// Otherwise, list all apps
			apps, err := client.GetApps()
			if err != nil {
				return fmt.Errorf("failed to list apps: %w", err)
			}

			// Print output
			formatted, err := output.FormatApps(apps, outputFormat)
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
	CommonFlagSet()(get)
	
	// Add command-specific flags
	get.Flags().StringP(constants.IDFlagName, constants.IDFlagShort, "", constants.IDFlagHelp)
}