package commands

import (
	"fmt"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.1"
	commit  = "dev"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     constants.RootCmdUse,
	Short:   "A command line interface wrapper for Humanitec platform",
	Long:    `A command line interface wrapper for Humanitec platform that provides basic CRUD operations for managing resources.`,
	Version: fmt.Sprintf("%s", version),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	// Add get command to root
	rootCmd.AddCommand(getCmd)
	return rootCmd.Execute()
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().StringP(constants.ConfigFlagName, constants.ConfigFlagShort, constants.DefaultConfigFile, "config file (default is $HOME/.humctl-wrapper.yaml)")
	rootCmd.PersistentFlags().BoolP(constants.VersionFlagName, constants.VersionFlagShort, false, "Print the version number")
} 