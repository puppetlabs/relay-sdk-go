package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/puppetlabs/nebula-sdk/pkg/container/def"
	"github.com/puppetlabs/nebula-sdk/pkg/container/defwalker"
	"github.com/puppetlabs/nebula-sdk/pkg/container/generator"
	"github.com/spf13/cobra"
)

var (
	verbose        bool
	write          bool
	recursive      bool
	scriptFilename string
	repoNameBase   string
)

func main() {
	c := &cobra.Command{
		Use:   os.Args[0],
		Short: "Generate Dockerfiles for Nebula containers from YAML configuration",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.OutOrStdout(), args)
		},
		SilenceUsage: true,
	}

	c.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print each file written")
	c.Flags().BoolVarP(&write, "write", "w", false, "Whether to persist the generated files")
	c.Flags().BoolVarP(&recursive, "recursive", "r", false, "Whether to recurse directory arguments")
	c.Flags().StringVar(&scriptFilename, "script-filename", generator.DefaultScriptFilename, "The file name to use when generating the build script")
	c.Flags().StringVarP(&repoNameBase, "repo-name-base", "n", generator.DefaultRepoNameBase, "The base Docker repository name to use")

	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(w io.Writer, args []string) error {
	if len(args) == 0 {
		args = []string{"."}
	}

	var containers []*def.ResolvedContainer

	for _, arg := range args {
		if !recursive {
			c, err := def.NewFromFilePath(arg)
			if err != nil {
				return err
			}

			containers = append(containers, c)
		} else {
			cs, err := defwalker.Walk(arg)
			if err != nil {
				return err
			}

			containers = append(containers, cs...)
		}
	}

	var fs []*generator.File
	for _, c := range containers {
		g := generator.New(
			c.Container,
			generator.WithFilesRelativeTo(c.FileRef),
			generator.WithScriptFilename(scriptFilename),
			generator.WithRepoNameBase(repoNameBase),
		)

		gfs, err := g.Files()
		if err != nil {
			return err
		}

		fs = append(fs, gfs...)
	}

	for _, f := range fs {
		dest := filepath.FromSlash(f.Ref.Name())

		prev, err := ioutil.ReadFile(dest)
		if os.IsNotExist(err) {
			// This is okay since we'll create the file anyway.
		} else if err != nil {
			return err
		}

		if !write {
			diff := difflib.UnifiedDiff{
				A:        splitLinesUnlessEmpty(string(prev)),
				B:        splitLinesUnlessEmpty(f.Content),
				FromFile: "a" + string(filepath.Separator) + dest,
				ToFile:   "b" + string(filepath.Separator) + dest,
				Context:  3,
			}
			if err := difflib.WriteUnifiedDiff(w, diff); err != nil {
				io.WriteString(w, formatError(fmt.Sprintf("Unable to generate changes for file %q", dest), err)+"\n")
				continue
			}
		} else if string(prev) != f.Content {
			if err := ioutil.WriteFile(dest, []byte(f.Content), f.Mode); err != nil {
				io.WriteString(w, formatError(fmt.Sprintf("Unable to generate file %q", dest), err)+"\n")
				continue
			}

			if err := os.Chmod(dest, f.Mode); err != nil {
				io.WriteString(w, formatError(fmt.Sprintf("Unable to set permissions for file %q", dest), err)+"\n")
				continue
			}

			if verbose {
				io.WriteString(w, dest+"\n")
			}
		}
	}

	return nil
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
