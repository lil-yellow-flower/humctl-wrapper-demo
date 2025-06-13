package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	commit  = "dev"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "humctl-wrapper-demo",
	Short: "A CLI wrapper for Humanitec platform",
	Long: `A command line interface wrapper for Humanitec platform that provides
basic CRUD operations for managing resources.`,
	Version: fmt.Sprintf("%s", version),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.humctl-wrapper.yaml)")
	
	// Add version flag
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")
	rootCmd.Flags().Lookup("version").NoOptDefVal = "true"
} 