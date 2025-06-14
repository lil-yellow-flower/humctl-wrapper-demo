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
)

func init() {
	rootCmd.AddCommand(getCmd)

	// Add output format flag
	getCmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, "", constants.OutputFlagHelp)
}