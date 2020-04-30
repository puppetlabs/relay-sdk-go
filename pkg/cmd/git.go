package cmd

import (
	"github.com/puppetlabs/nebula-sdk/pkg/task"
	"github.com/puppetlabs/nebula-sdk/pkg/taskutil"
	"github.com/spf13/cobra"
)

func NewGitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "git",
		Short:                 "Manage git data",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewGitCloneCommand())

	return cmd
}

func NewGitCloneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "clone",
		Short:                 "Clone git repository",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			directory, _ := cmd.Flags().GetString("directory")
			revision, _ := cmd.Flags().GetString("revision")

			planOpts := taskutil.DefaultPlanOptions{SpecURL: taskutil.MetadataSpecURL()}
			task := task.NewTaskInterface(planOpts)
			return task.CloneRepository(revision, directory)
		},
	}

	cmd.Flags().StringP("revision", "r", "", "git revision")
	cmd.Flags().StringP("directory", "d", "", "git clone output directory")

	return cmd
}
