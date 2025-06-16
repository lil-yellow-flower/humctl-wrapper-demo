package apps

import (
	"github.com/spf13/cobra"
)

// GetCommand returns the command for getting apps
func GetCommand() *cobra.Command {
	return get
}

// CreateCommand returns the command for creating apps
func CreateCommand() *cobra.Command {
	return create
}

// UpdateCommand returns the command for updating apps
func UpdateCommand() *cobra.Command {
	return update
}

// DeleteCommand returns the command for deleting apps
func DeleteCommand() *cobra.Command {
	return delete
} 