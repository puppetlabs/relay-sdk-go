package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/puppetlabs/nebula-sdk/pkg/container/def"
	"github.com/puppetlabs/nebula-sdk/pkg/container/defwalker"
)

func containers(paths []string, recursive bool) ([]*def.ResolvedContainer, error) {
	var containers []*def.ResolvedContainer

	for _, path := range paths {
		if !recursive {
			c, err := def.NewFromFilePath(path)
			if err != nil {
				return nil, err
			}

			containers = append(containers, c)
		} else {
			cs, err := defwalker.Walk(path)
			if err != nil {
				return nil, err
			}

			containers = append(containers, cs...)
		}
	}

	return containers, nil
}

func splitLinesUnlessEmpty(in string) []string {
	if in == "" {
		return nil
	}

	return difflib.SplitLines(in)
}

func formatError(explanation string, err error) string {
	s := fmt.Sprintf("# %s: %+v", explanation, err)
	return strings.ReplaceAll(s, "\n", "\n#")
}

var wd, _ = os.Getwd()

func relativePath(fr *def.FileRef) string {
	if !fr.Local() {
		return fr.String()
	}

	p := filepath.FromSlash(fr.String())
	if rp, err := filepath.Rel(wd, p); err == nil {
		p = rp
	}

	return p
}

func writeTable(w io.Writer, fn func(t table.Writer)) {
	t := table.NewWriter()
	t.SetOutputMirror(w)
	t.SetStyle(table.Style{
		Box:     table.StyleBoxDefault,
		Color:   table.ColorOptionsDefault,
		Format:  table.FormatOptionsDefault,
		Options: table.OptionsNoBordersAndSeparators,
		Title:   table.TitleOptionsDefault,
	})

	fn(t)

	t.Render()
}
