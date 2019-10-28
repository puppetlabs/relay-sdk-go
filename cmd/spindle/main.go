package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:          os.Args[0],
		Short:        "Generate Dockerfiles for Nebula containers from YAML configuration",
		SilenceUsage: true,
	}

	cmd.AddCommand(NewGenerateCommand())
	cmd.AddCommand(NewListCommand())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
