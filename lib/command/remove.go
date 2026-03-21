package command

import (
	"takyon/lib/container"
	"takyon/lib/container/cutils"
	"takyon/lib/ui"

	"github.com/spf13/cobra"
)

var force bool

var RemoveCMD = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a disk image",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			ui.Warn("The container's name is required")
			return
		}

		name := args[0]

		if !cutils.ImageExist(name) {
			ui.Error("No container of this name exist")
			return
		}

		if cutils.IsMounted(name) {

			if force {
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

		img_path := cutils.GetImagePath(name)
		ui.Step("Removing disk image %s", img_path)
		err := cutils.Run("rm", img_path)

		if err != nil {
			ui.Fatal(err.Error())
			return
		}

		ui.Success("Successfully removed container %s", name)
	},
}

func init() {
	RemoveCMD.Flags().BoolVarP(&force, "force", "f", false, "Automaaticaly Umount the image")
}
