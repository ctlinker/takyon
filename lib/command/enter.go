package command

import (
	"takyon/lib/container"
	"takyon/lib/ui"

	"github.com/spf13/cobra"
)

var shell string
var user string

var EnterCMD = &cobra.Command{
	Use:   "enter [name]",
	Short: "Chroot inside a container image",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			ui.Warn("The container's name is required")
			return
		}

		name := args[0]

		err := container.EnterContainer(container.EnterContainerOption{
			Name:  name,
			User:  user,
			Shell: shell,
		})

		if err != nil {
			ui.Fatal(err.Error())
			return
		}
	},
}

func init() {
	EnterCMD.Flags().StringVarP(&shell, "shell", "s", "/bin/bash", "Shell")
	EnterCMD.Flags().StringVarP(&user, "user", "u", "root", "Username")
}
