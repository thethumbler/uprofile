package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/thethumbler/uprofile/internal"
)

var JumpCommand = &cobra.Command{
	Use:   "jump <profile>",
	Short: "jump to profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		context := internal.GetContext()
		profile := args[0]

		profileDir := fmt.Sprintf("%s/%s", context.ProfilesDir, profile)
		homeDir := fmt.Sprintf("%s/merged", profileDir)

		if info, err := os.Stat(profileDir); err != nil {
			if err.Error() == "no such file or directory" {
				fmt.Fprintf(os.Stderr, "no such profile: %s\n", profile)
			} else {
				fmt.Fprintf(os.Stderr, "unable to access profile '%s', '%s': %s\n", profile, profileDir, err.Error())
			}

			os.Exit(1)
		} else if !info.IsDir() {
			fmt.Fprintf(os.Stderr, "profile '%s' is not created correctly, '%s' is not a directory.\n", profile, profileDir)
			os.Exit(1)
		}

		env := append(os.Environ(),
			fmt.Sprintf("HOME=%s", homeDir),
			fmt.Sprintf("PS1=[\\u@%s \\W]\\$ ", profile),
		)

		binaryPath := "/usr/bin/unshare"
		if err := syscall.Exec(binaryPath, []string{"unshare", "-w", homeDir}, env); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run 'unshare': %s", err.Error())
			os.Exit(1)
		}
	},
}
