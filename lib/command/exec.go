package command

import (
	"takyon/lib/container"
	"takyon/lib/ui"

	"github.com/spf13/cobra"
)

var script string
var as_user string

var ExecCMD = &cobra.Command{
	Use:   "exec [name]",
	Short: "Exec a command inside a container",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			ui.Warn("The container's name is required")
			return
		}

		name := args[0]

		err := container.ExecInContainer(container.ExecInContainerOption{
			Name:   name,
			User:   as_user,
			Script: script,
		})

		if err != nil {
			ui.Fatal(err.Error())
			return
		}
	},
}

func init() {
	ExecCMD.Flags().StringVarP(&script, "script", "s", "", "Command or script content")
	ExecCMD.Flags().StringVarP(&as_user, "user", "u", "root", "Username")
}
