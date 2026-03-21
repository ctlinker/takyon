package command

import (
	"takyon/lib/container"
	"takyon/lib/ui"

	"github.com/spf13/cobra"
)

var UmountCMD = &cobra.Command{
	Use:   "umount [name]",
	Short: "Umount a new environment",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			ui.Warn("The container's name is required")
			return
		}

		name := args[0]

		err := container.UmountDiskImage(name)
		if err != nil {
			ui.Fatal(err.Error())
		}
	},
}
