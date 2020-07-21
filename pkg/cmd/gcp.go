package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/puppetlabs/relay-sdk-go/pkg/task"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"github.com/spf13/cobra"
)

func NewGCPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "gcp",
		Short:                 "Manage Google Cloud Platform access",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewGCPConfigCommand())
	cmd.AddCommand(NewGCPEnvCommand())

	return cmd
}

func NewGCPConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "config",
		Short:                 "Create a GCP configuration suitable for using with a GCP CLI or SDK",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			directory, _ := cmd.Flags().GetString("directory")

			u, err := taskutil.MetadataSpecURL()
			if err != nil {
				return err
			}
			planOpts := taskutil.DefaultPlanOptions{SpecURL: u}
			task := task.NewTaskInterface(planOpts)
			return task.ProcessGCP(directory)
		},
	}

	cmd.Flags().StringP("directory", "d", "", "configuration output directory")

	return cmd
}

func NewGCPEnvCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "env",
		Short:                 "Create a POSIX-compatible script that can be sourced to configure the GCP CLI",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			directory, _ := cmd.Flags().GetString("directory")

			fmt.Fprintf(
				cmd.OutOrStdout(),
				`export GOOGLE_APPLICATION_CREDENTIALS=%s`,
				quoteShell(filepath.Join(directory, "credentials.json")),
			)
			return nil
		},
	}

	cmd.Flags().StringP("directory", "d", "", "configuration output directory")

	return cmd
}
