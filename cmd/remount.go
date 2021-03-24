package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/thethumbler/uprofile/internal"
)

var RemountCommand = &cobra.Command{
	Use:   "remount <profile>",
	Short: "remount a profile to reflect changes done in the original home directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		context := internal.GetContext()
		profile := args[0]

		if err := internal.UmountProfile(&context, profile); err != nil {
			log.Fatal(err)
		}

		if err := internal.MountProfile(&context, profile); err != nil {
			log.Fatal(err)
		}
	},
}
