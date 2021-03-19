package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thethumbler/uprofile/internal"
)

var DeleteCommand = &cobra.Command{
	Use:   "delete <profile>",
	Short: "delete profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		context := internal.GetContext()
		profile := args[0]

		profileDir := fmt.Sprintf("%s/%s", context.ProfilesDir, profile)
		os.RemoveAll(profileDir)
	},
}
