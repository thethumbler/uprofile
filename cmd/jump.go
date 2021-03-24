package cmd

import (
	"log"

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

		if err := internal.Jump(&context, profile); err != nil {
			log.Fatal(err)
		}
	},
}
