package main

import (
	"io"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

type listCommand struct {
	quiet     bool
	recursive bool
}

func (lc *listCommand) run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		args = []string{"."}
	}

	cs, err := containers(args, lc.recursive)
	if err != nil {
		return err
	}

	if lc.quiet {
		for _, c := range cs {
			io.WriteString(cmd.OutOrStdout(), relativePath(c.FileRef)+"\n")
		}
	} else {
		writeTable(cmd.OutOrStdout(), func(t table.Writer) {
			t.AppendHeader(table.Row{"ID", "Name", "SDK Version", "Path"})

			for _, c := range cs {
				t.AppendRow(table.Row{
					c.Container.ID[:12],
					c.Container.Name,
					c.Container.SDKVersion,
					relativePath(c.FileRef),
				})
			}
		})
	}

	return nil
}

func NewListCommand() *cobra.Command {
	lc := &listCommand{}

	cmd := &cobra.Command{
		Use:     "list [flags] <paths...>",
		Short:   "List the containers detected in the given files or directories",
		Aliases: []string{"ls"},
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return lc.run(cmd, args)
		},
	}

	cmd.Flags().BoolVarP(&lc.quiet, "quiet", "q", false, "only print the paths to the container descriptor files instead of a table")
	cmd.Flags().BoolVarP(&lc.recursive, "recursive", "r", false, "whether to recurse subdirectories of arguments")

	return cmd
}
