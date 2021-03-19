package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jsgilmore/mount"
	"github.com/spf13/cobra"
	"github.com/thethumbler/uprofile/internal"
)

var listAll bool = false

var ListCommand = &cobra.Command{
	Use:   "list",
	Short: "list profiles",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		context := internal.GetContext()

		if listAll {
			profilesList, _ := ioutil.ReadDir(context.ProfilesDir)
			for _, entry := range profilesList {
				fmt.Println(entry.Name())
			}
		} else {
			mountedProfiles, _ := mount.Mounts()
			for _, mount := range mountedProfiles {
				if mount.Filesystem == "fuse.fuse-overlayfs" {
					var match string
					n, _ := fmt.Sscanf(mount.Path, fmt.Sprintf("%s/%%s", context.ProfilesDir), &match)
					if n == 1 {
						profile := strings.Split(match, "/")[0]
						fmt.Println(profile)
					}
				}
			}
		}
	},
}

func init() {
	ListCommand.PersistentFlags().BoolVarP(&listAll, "all", "a", false, "list all profiles")
}
