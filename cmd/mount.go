package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/thethumbler/uprofile/internal"
)

var MountCommand = &cobra.Command{
	Use:   "mount <profile>",
	Short: "mount profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		context := internal.GetContext()
		profile := args[0]

		if err := internal.MountProfile(&context, profile); err != nil {
			log.Fatal(err)
		}
	},
}
