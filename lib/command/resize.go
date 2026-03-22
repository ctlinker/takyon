package command

import (
	"strconv"
	"takyon/lib/container"
	"takyon/lib/container/cutils"
	"takyon/lib/ui"

	"github.com/spf13/cobra"
)

var forceUmount bool

var ResizeCMD = &cobra.Command{
	Use:   "resize [name] [size]",
	Short: "Resize a container image",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			ui.Warn("The container's name and the size are required")
			return
		}

		name := args[0]
		size_str := args[1]

		unit := size_str[len(size_str)-1:]
		size, str_err := strconv.Atoi(size_str[:len(size_str)-1])

		if str_err != nil {
			ui.Warn("Invalid container size, expecting an integer")
			ui.Error(ui.AbortErr.Error())
			return
		}

		if cutils.IsMounted(name) {

			if forceUmount {
				err := container.UmountDiskImage(name)
				if err != nil {
					ui.Fatal(err.Error())
					return
				}
			} else {
				ui.Warn("The Container %s is currently mounted, umount it or use -f flag", name)
				ui.Error(ui.AbortErr.Error())
				return
			}

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
	ResizeCMD.Flags().BoolVarP(&forceUmount, "forceUmount", "f", false, "Automaaticaly Umount the image")
}
