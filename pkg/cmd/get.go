package cmd

import (
	"github.com/puppetlabs/relay-sdk-go/pkg/task"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"github.com/spf13/cobra"
)

func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "get",
		Short:                 "Get specification data",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("path")

			u, err := taskutil.MetadataSpecURL()
			if err != nil {
				return err
			}
			planOpts := taskutil.DefaultPlanOptions{SpecURL: u}
			task := task.NewTaskInterface(planOpts)
			data, err := task.ReadData(path)

			if err != nil {
				return err
			}

			cmd.OutOrStdout().Write(data)

			return nil
		},
	}

	cmd.Flags().StringP("path", "p", "", "specification data path")

	return cmd
}
