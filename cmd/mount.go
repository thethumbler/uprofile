package cmd

import (
	"fmt"
	"os/exec"

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

		lowerDir := fmt.Sprintf("/home/%s", context.User)
		upperDir := fmt.Sprintf("%s/%s/upperdir", context.ProfilesDir, profile)
		workDir := fmt.Sprintf("%s/%s/workdir", context.ProfilesDir, profile)
		mergedDir := fmt.Sprintf("%s/%s/merged", context.ProfilesDir, profile)

		mountOptions := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDir, upperDir, workDir)

		exec.Command("fuse-overlayfs", "-o", mountOptions, mergedDir).Run()
	},
}
