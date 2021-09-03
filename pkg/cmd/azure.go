package cmd

import (
	"fmt"

	"github.com/puppetlabs/relay-sdk-go/pkg/task"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"github.com/spf13/cobra"
)

func NewAzureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "azure",
		Short:                 "Manage Azure access",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewAzureARMCommand())

	return cmd
}

func NewAzureARMCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "arm",
		Short:                 "Manage ARM access",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewAzureARMEnvCommand())

	return cmd
}

func NewAzureARMEnvCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "env",
		Short:                 "Create a POSIX-compatible script that can be sourced to configure Azure ARM",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			u, err := taskutil.MetadataSpecURL()
			if err != nil {
				return err
			}
			planOpts := taskutil.DefaultPlanOptions{SpecURL: u}
			task := task.NewTaskInterface(planOpts)
			env, err := task.GetAzureARMEnvironmentVariables()
			if err != nil {
				return err
			}

			export := ""
			for k, v := range env {
				export += fmt.Sprintf("export %s=%s\n", k, quoteShell(v))
			}

			fmt.Fprint(cmd.OutOrStdout(), export)

			return nil
		},
	}

	return cmd
}
