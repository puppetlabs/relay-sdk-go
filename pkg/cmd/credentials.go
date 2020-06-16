package cmd

import (
	"github.com/puppetlabs/relay-sdk-go/pkg/task"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"github.com/spf13/cobra"
)

func NewCredentialsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "credentials",
		Short:                 "Manage credentials configuration",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewCredentialsConfigCommand())

	return cmd
}

func NewCredentialsConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "config",
		Short:                 "Create credentials configuration",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			directory, _ := cmd.Flags().GetString("directory")

			u, err := taskutil.MetadataSpecURL()
			if err != nil {
				return err
			}
			planOpts := taskutil.DefaultPlanOptions{SpecURL: u}
			task := task.NewTaskInterface(planOpts)
			return task.ProcessCredentials(directory)
		},
	}

	cmd.Flags().StringP("directory", "d", "", "configuration output directory")

	return cmd
}
