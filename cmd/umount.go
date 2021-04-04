package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/thethumbler/uprofile/internal"
)

var unmountAll bool = false

var UmountCommand = &cobra.Command{
	Use:   "umount [profile]...",
	Short: "unmount profile",
	Run: func(cmd *cobra.Command, args []string) {
		context := internal.GetContext()
		unmountProfiles := args

		if unmountAll {
			unmountProfiles = internal.MountedProfiles(&context)
		}

		for _, profile := range unmountProfiles {
			if err := internal.UmountProfile(&context, profile); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	UmountCommand.PersistentFlags().BoolVarP(&unmountAll, "all", "a", false, "umount all mounted profiles")
}
