package main

import (
	"fmt"
	"os"

	"takyon/lib/command"
	"takyon/lib/env"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "takyon",
		Short: "Takyon container manager",
	}

	// attach subcommands
	rootCmd.AddCommand(command.CreateCMD)
	rootCmd.AddCommand(command.ListCMD)
	rootCmd.AddCommand(command.MountCMD)
	rootCmd.AddCommand(command.FlashCMD)
	rootCmd.AddCommand(command.EnterCMD)

	env.SetupDir()

	// run CLI
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
