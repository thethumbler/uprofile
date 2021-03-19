package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thethumbler/uprofile/internal"
)

var CreateCommand = &cobra.Command{
	Use:   "create <profile>",
	Short: "create a new profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		context := internal.GetContext()
		profile := args[0]

		dirsList := []string{
			fmt.Sprintf("%s/%s/upperdir", context.ProfilesDir, profile),
			fmt.Sprintf("%s/%s/workdir", context.ProfilesDir, profile),
			fmt.Sprintf("%s/%s/merged", context.ProfilesDir, profile),
		}

		for _, path := range dirsList {
			os.MkdirAll(path, 0700)
		}
	},
}
