package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/thethumbler/uprofile/internal"
)

var UmountCommand = &cobra.Command{
	Use:   "umount <profile>",
	Short: "unmount profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		context := internal.GetContext()
		profile := args[0]

		mountpoint := fmt.Sprintf("%s/%s/merged", context.ProfilesDir, profile)
		exec.Command("umount", mountpoint).Run()
	},
}
