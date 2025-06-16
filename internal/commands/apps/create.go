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
	// Subcommand for creating apps
	create = &cobra.Command{
		Use:   constants.AppCmdUse,
		Short: constants.AppCmdShort,
		Long: `Create a new application in the organization.
The application requires both a name (human-friendly display name) and an ID (unique identifier).
The ID must match the pattern: ^[a-z0-9](?:-?[a-z0-9]+)+$`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flags
			name, err := cmd.Flags().GetString(constants.NameFlagName)
			if err != nil {
				return fmt.Errorf("failed to get name flag: %w", err)
			}

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

			skipEnvCreation, err := cmd.Flags().GetBool(constants.SkipEnvCreationFlagName)
			if err != nil {
				return fmt.Errorf("failed to get skip environment creation flag: %w", err)
			}

			// Get organization ID from config
			org := config.GetConfig().HumanitecOrg
			token := config.GetConfig().HumanitecToken

			// Create Humanitec client
			client := humanitec.NewClient(token, org)

			// Create app
			app, err := client.CreateApp(id, name, skipEnvCreation)
			if err != nil {
				return fmt.Errorf("failed to create app: %w", err)
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
	CommonFlagSet()(create)
	
	// Add command-specific flags
	create.Flags().StringP(constants.NameFlagName, constants.NameFlagShort, "", constants.NameFlagHelp)
	create.Flags().StringP(constants.IDFlagName, constants.IDFlagShort, "", constants.IDFlagHelp)
	create.Flags().BoolP(constants.SkipEnvCreationFlagName, constants.SkipEnvCreationFlagShort, false, constants.SkipEnvCreationFlagHelp)
	create.MarkFlagRequired(constants.NameFlagName)
	create.MarkFlagRequired(constants.IDFlagName)
} 