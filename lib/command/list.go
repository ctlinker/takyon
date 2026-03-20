package command

import (
	"takyon/lib/container"

	"github.com/spf13/cobra"
)

var ListCMD = &cobra.Command{
	Use:   "list",
	Short: "List the availables container images",
	Run: func(cmd *cobra.Command, args []string) {
		container.ListContainers()
	},
}
