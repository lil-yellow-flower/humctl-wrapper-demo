package apps

import (
	"fmt"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/config"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/output"
	"github.com/spf13/cobra"
)

var delete = &cobra.Command{
	Use:   constants.AppCmdUse,
	Short: constants.AppCmdShort,
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

		// Delete app
		if err := client.DeleteApp(id); err != nil {
			return fmt.Errorf("failed to delete app: %w", err)
		}

		// Print output
		formatted, err := output.FormatMessage("Application successfully deleted", outputFormat)
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}
		fmt.Fprint(cmd.OutOrStdout(), formatted)

		return nil
	},
}

func init() {
	// Add common flags
	CommonFlagSet()(delete)
	
	// Add command-specific flags
	delete.Flags().StringP(constants.IDFlagName, constants.IDFlagShort, "", constants.IDFlagHelp)
	delete.MarkFlagRequired(constants.IDFlagName)
} 