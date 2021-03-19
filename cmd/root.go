package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "uprofile",
	Short: "Manage user profiles",
	Long:  "Manage multiple profiles for a single user using overlay mounts for each profile",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello world")
	},
}

func Execute() {
	rootCommand.AddCommand(CreateCommand)
	rootCommand.AddCommand(DeleteCommand)
	rootCommand.AddCommand(ListCommand)
	rootCommand.AddCommand(MountCommand)
	rootCommand.AddCommand(UmountCommand)
	rootCommand.AddCommand(JumpCommand)

	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
