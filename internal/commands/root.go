package commands

import (
	"fmt"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/config"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/commands/apps"
	"github.com/spf13/cobra"
)

var (
	// These variables are set during build time using ldflags
	version = "0.5.0"
	commit  = "unknown"
	date    = "unknown"
)

var RootCmd = &cobra.Command{
	Use:     "humctl-wrapper",
	Short:   "A wrapper for the Humanitec CLI",
	Long:    `A wrapper for the Humanitec CLI that provides additional functionality and a more user-friendly interface.`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Initialize("config.yaml"); err != nil {
			return fmt.Errorf("error loading config: %v", err)
		}

		cfg := config.GetConfig()
		if cfg.HumanitecToken == "" {
			return fmt.Errorf(constants.ErrMissingToken)
		}
		if cfg.HumanitecOrg == "" {
			return fmt.Errorf(constants.ErrMissingOrg)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	// Add get command
	getCmd := &cobra.Command{
		Use:   constants.GetCmdUse,
		Short: constants.GetCmdShort,
	}
	RootCmd.AddCommand(getCmd)
	
	// Add create command
	createCmd := &cobra.Command{
		Use:   constants.CreateCmdUse,
		Short: constants.CreateCmdShort,
	}
	RootCmd.AddCommand(createCmd)

	// Add update command
	updateCmd := &cobra.Command{
		Use:   constants.UpdateCmdUse,
		Short: constants.UpdateCmdShort,
	}
	RootCmd.AddCommand(updateCmd)

	// Add delete command
	deleteCmd := &cobra.Command{
		Use:   constants.DeleteCmdUse,
		Short: constants.DeleteCmdShort,
	}
	RootCmd.AddCommand(deleteCmd)

	// Add apps as subcommand of each verb
	getCmd.AddCommand(apps.GetCommand())
	createCmd.AddCommand(apps.CreateCommand())
	updateCmd.AddCommand(apps.UpdateCommand())
	deleteCmd.AddCommand(apps.DeleteCommand())

	return RootCmd.Execute()
}

func init() {
	RootCmd.PersistentFlags().BoolP(constants.VersionFlagName, constants.VersionFlagShort, false, "Print the version number")
} 