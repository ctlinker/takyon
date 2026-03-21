package command

import (
	"strconv"
	"takyon/lib/container"
	"takyon/lib/ui"

	"github.com/spf13/cobra"
)

var unit string

var ResizeCMD = &cobra.Command{
	Use:   "resize [name] [size]",
	Short: "Resize a container image",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			ui.Warn("The container's name and the size are required")
			return
		}

		name := args[0]
		size, str_err := strconv.Atoi(args[1])

		if str_err != nil {
			ui.Warn("Invalid container size, expecting an integer")
			return
		}

		err := container.ResizeDiskImage(container.ResizeDiskImageOption{
			Name: name,
			Unit: unit,
			Size: size,
		})

		if err != nil {
			ui.Fatal(err.Error())
			return
		}
	},
}

func init() {
	ResizeCMD.Flags().StringVarP(&unit, "unit", "u", "M", "Size unit, defaut M for mega")
}
