package command

import (
	"takyon/lib/container"
	"takyon/lib/ui"

	"github.com/spf13/cobra"
)

var size int
var format string

var CreateCMD = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new container image",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			ui.Warn("The container's name is required")
			return
		}

		name := args[0]

		err := container.CreateDiskImage(container.CreateDiskImageOption{
			Name:   name,
			Format: format,
			Size:   size,
		})

		if err != nil {
			ui.Fatal(err.Error())
			return
		}
	},
}

func init() {
	CreateCMD.Flags().IntVarP(&size, "size", "s", 2048, "Size in MB")
	CreateCMD.Flags().StringVarP(&format, "format", "f", "ext4", "Filesystem format")
}
