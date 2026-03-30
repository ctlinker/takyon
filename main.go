package main

import (
	"fmt"
	"os"

	"takyon/library/env"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "takyon command [-S env-store] [-M env-mount-dir]",
		Short:   "Takyon container manager",
		Version: "2.0.0",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := env.SetupEnvDirectories(); err != nil {
				return
			}
		},
	}

	// run CLI
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
