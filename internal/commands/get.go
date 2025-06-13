package commands

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/output"
	"github.com/spf13/cobra"
)

var (
	client humanitec.Client
)

// SetClient sets the Humanitec client
func SetClient(c humanitec.Client) {
	client = c
}

var getCmd = &cobra.Command{
	Use:   constants.GetCmdUse,
	Short: "Get resources from Humanitec platform",
	Long:  `Get resources from Humanitec platform. For example: humctl-wrapper get apps`,
}

var getAppsCmd = &cobra.Command{
	Use:   constants.GetAppsCmdUse,
	Short: "Get applications from Humanitec platform",
	Long:  `Get applications from Humanitec platform. For example: humctl-wrapper get apps`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load .env file
		if err := godotenv.Load(); err != nil {
			// Don't return error if .env file doesn't exist
			fmt.Println("Warning: .env file not found")
		}

		// Get format from flag
		formatStr, err := cmd.Flags().GetString(constants.OutputFlagName)
		if err != nil {
			return fmt.Errorf(constants.ErrInvalidOutputFormat, err)
		}

		// Validate format
		format, err := output.ValidateFormat(formatStr)
		if err != nil {
			return fmt.Errorf(constants.ErrInvalidOutputFormat, err)
		}

		// Get org from flag or environment
		org, err := cmd.Flags().GetString(constants.OrgFlagName)
		if err != nil {
			return fmt.Errorf(constants.ErrInvalidOrgFlag, err)
		}

		// If org not provided, use environment variable
		if org == "" {
			org = os.Getenv(constants.EnvHumanitecOrg)
			if org == "" {
				return fmt.Errorf(constants.ErrMissingOrg)
			}
		}

		// Get token from environment
		token := os.Getenv(constants.EnvHumanitecToken)
		if token == "" {
			return fmt.Errorf(constants.ErrMissingToken)
		}

		// Initialize client
		c := humanitec.NewClient(token, org)

		// Get applications
		apps, err := c.ListApps()
		if err != nil {
			return fmt.Errorf(constants.ErrGetApps, err)
		}

		// Format and print output
		formatted, err := output.FormatApps(apps, format)
		if err != nil {
			return fmt.Errorf(constants.ErrFormatOutput, err)
		}
		fmt.Println(formatted)

		return nil
	},
}

func init() {
	// Add get apps command to get command
	getCmd.AddCommand(getAppsCmd)

	// Add flags to get apps command
	getAppsCmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, constants.DefaultOutputFormat, "Output format (table|json|yaml)")
	getAppsCmd.Flags().String(constants.OrgFlagName, "", fmt.Sprintf("Humanitec organization ID (defaults to %s environment variable)", constants.EnvHumanitecOrg))
} 