// go:generate go run generate_docs.go

package cmd

import (
	"bytes"

	"github.com/spf13/cobra"
)

// NewDocCommand adds the 'doc' subcommand under the root 'ni'
func NewDocCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doc",
		Short: "build command documentation",
		Args:  cobra.NoArgs,
		RunE:  genDocs,
	}

	return cmd
}

func genDocs(cmd *cobra.Command, args []string) error {
	buf := new(bytes.Buffer)
	rootcmd, err := NewRootCommand()

	if err != nil {
		return err
	}

	rootcmd.InitDefaultHelpCmd()
	rootcmd.InitDefaultHelpFlag()

	name := rootcmd.CommandPath()

	short := rootcmd.Short
	long := rootcmd.Long

	buf.WriteString("## " + name + "\n\n" + short + "\n\n")
	buf.WriteString("### Synopsis\n\n" + long + "\n\n")
	buf.WriteString("### Subcommand Usage\n\n")

	children := rootcmd.Commands()

	if err := genChildMarkdown(children, buf); err != nil {
		return err
	}

	flags := rootcmd.PersistentFlags()
	if flags.HasAvailableFlags() {
		buf.WriteString("### Global flags\n```\n")
		buf.WriteString(flags.FlagUsages() + "\n```\n")
	}
	rootcmd.OutOrStdout().Write(buf.Bytes())

	return nil

}

// For brevity, we only want to generate docs for 'leaf' commands, i.e.
// only "ni output get", not "ni output"
func genChildMarkdown(children []*cobra.Command, buf *bytes.Buffer) error {

	for _, child := range children {
		if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
			continue
		}
		if child.Runnable() {
			usage := child.UseLine()
			buf.WriteString("**`" + usage + "`** -- " + child.Short + "\n")
			long := child.Long
			if len(long) > 0 {
				buf.WriteString("  " + child.Long + "\n")
			}
			flags := child.NonInheritedFlags()
			if flags.HasAvailableFlags() {
				buf.WriteString("```\n")
				flags.SetOutput(buf)
				flags.PrintDefaults()
				buf.WriteString("```\n")
			}
			buf.WriteString("\n")
		}
		// Because commands can be nested arbitrarily deep, this recurses into
		// the current command's children and generates their docs too
		grandchildren := child.Commands()
		genChildMarkdown(grandchildren, buf)

	}

	return nil

}
