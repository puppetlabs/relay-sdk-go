package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// NewDocCommand adds the 'doc' subcommand under the root 'ni'
func NewDocCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doc",
		Short: "build command documentation",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newGenerateCommand())

	return cmd
}

func newGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate documentation",
		Args:  cobra.NoArgs,
		RunE:  genDocs,
	}

	cmd.Flags().StringP("target", "t", "./docs/", "target directory for output")

	return cmd
}

func genDocs(cmd *cobra.Command, args []string) error {
	targetDir, nil := cmd.Flags().GetString("target")

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		cmd.Printf(`Failed to create target directory %s: %s`, targetDir, err.Error())
		return err
	}

	rootcmd, nil := NewRootCommand()

	doc.GenMarkdownTree(rootcmd, targetDir)

	return nil

}
