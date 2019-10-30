package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/puppetlabs/nebula-sdk/pkg/container/def"
	"github.com/spf13/cobra"
)

type describeTemplateCommand struct {
	from string
}

func (dtc *describeTemplateCommand) run(cmd *cobra.Command, arg string) error {
	var opts []def.FileRefOption
	switch dtc.from {
	case "":
	case "sdk":
		opts = append(opts, def.WithFileRefResolver(def.SDKResolver))
	default:
		return fmt.Errorf("unknown file source %q", dtc.from)
	}

	tpl, err := def.NewTemplateFromFileRef(def.NewFileRef(arg, opts...))
	if err != nil {
		return err
	}

	io.WriteString(cmd.OutOrStdout(), "Images:\n")
	writeTable(cmd.OutOrStdout(), func(t table.Writer) {
		t.AppendHeader(table.Row{"Name", "Template", "Dependencies"})

		for name, image := range tpl.Images {
			t.AppendRow(table.Row{name, image.TemplateName, strings.Join(image.DependsOn, ", ")})
		}
	})

	io.WriteString(cmd.OutOrStdout(), "\nSettings:\n")
	writeTable(cmd.OutOrStdout(), func(t table.Writer) {
		t.AppendHeader(table.Row{"Name", "Description", "Default Value"})

		for name, setting := range tpl.Settings {
			defaultValue, _ := json.Marshal(setting.Value)
			t.AppendRow(table.Row{name, setting.Description, string(defaultValue)})
		}
	})

	return nil
}

func NewDescribeTemplateCommand() *cobra.Command {
	dtc := &describeTemplateCommand{}

	cmd := &cobra.Command{
		Use:     "template [flags] <path>",
		Short:   "Describe the given template",
		Aliases: []string{"tpl", "tmpl"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dtc.run(cmd, args[0])
		},
	}

	cmd.Flags().StringVarP(&dtc.from, "from", "f", "", "the source to load the template from if not the current file system")

	return cmd
}

func NewDescribeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Describe Spindle resources",
		Aliases: []string{"desc"},
	}

	cmd.AddCommand(NewDescribeTemplateCommand())

	return cmd
}
