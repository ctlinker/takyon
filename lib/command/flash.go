package command

import (
	"takyon/lib/container"
	"takyon/lib/ui"

	"github.com/spf13/cobra"
)

var distro string

// FlashCMD bootstraps a container with a minimal distro
var FlashCMD = &cobra.Command{
	Use:   "flash [container-name]",
	Short: "Flash a container with a minimal distro",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if err := container.FlashContainer(name, distro); err != nil {
			ui.Error("Failed to flash container %s: %v", name, err)
			return
		}
	},
}

func init() {
	// default distro = debian
	FlashCMD.Flags().StringVarP(&distro, "distro", "d", "debian", "Distro to flash (debian, ubuntu, kali, arch)")
}
