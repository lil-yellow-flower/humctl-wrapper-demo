package apps

import (
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/constants"
	"github.com/spf13/cobra"
)

// CommonFlagSet returns a function that adds common flags to a command
func CommonFlagSet() func(*cobra.Command) {
	return func(cmd *cobra.Command) {
		cmd.Flags().StringP(constants.OutputFlagName, constants.OutputFlagShort, constants.DefaultOutputFormat, constants.OutputFlagHelp)
		cmd.Flags().StringP(constants.OrgFlagName, constants.OrgFlagShort, "", constants.OrgFlagHelp)
	}
} 