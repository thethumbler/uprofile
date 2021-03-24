package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/thethumbler/uprofile/internal"
)

var remountAll bool = false

var RemountCommand = &cobra.Command{
	Use:   "remount [profile]...",
	Short: "remount a profile to reflect changes done in the original home directory",
	Run: func(cmd *cobra.Command, args []string) {
		context := internal.GetContext()

		remountProfiles := args

		if remountAll {
			remountProfiles = internal.MountedProfiles(&context)
		}

		for _, profile := range remountProfiles {
			if err := internal.UmountProfile(&context, profile); err != nil {
				log.Fatal(err)
			}

			if err := internal.MountProfile(&context, profile); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	RemountCommand.PersistentFlags().BoolVarP(&remountAll, "all", "a", false, "remount all mounted profiles")
}
