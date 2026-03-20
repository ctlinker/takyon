package command

import (
	"takyon/lib/container"
	"takyon/lib/ui"

	"github.com/spf13/cobra"
)

var MountCMD = &cobra.Command{
	Use:   "mount [name]",
	Short: "Mount a new environment",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			ui.Warn("The container's name is required")
			return
		}

		name := args[0]

		err := container.MountDiskImage(name)
		if err != nil {
			ui.Fatal(err.Error())
		}
	},
}
