package cmd

import (
	"github.com/puppetlabs/nebula-sdk/pkg/task"
	"github.com/puppetlabs/nebula-sdk/pkg/taskutil"
	"github.com/spf13/cobra"
)

func NewClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "cluster",
		Short:                 "Manage cluster configuration",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewClusterConfigCommand())

	return cmd
}

func NewClusterConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "config",
		Short:                 "Create cluster config",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			directory, _ := cmd.Flags().GetString("directory")

			u, err := taskutil.MetadataSpecURL()
			if err != nil {
				return err
			}
			planOpts := taskutil.DefaultPlanOptions{SpecURL: u}
			task := task.NewTaskInterface(planOpts)
			return task.ProcessClusters(directory)
		},
	}

	cmd.Flags().StringP("directory", "d", "", "configuration output directory")

	return cmd
}
