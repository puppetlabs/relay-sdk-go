package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/puppetlabs/nebula-sdk/pkg/container/generator"
	"github.com/spf13/cobra"
)

type generateCommand struct {
	verbose        bool
	write          bool
	recursive      bool
	scriptFilename string
	repoNameBase   string
}

func (gc *generateCommand) run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		args = []string{"."}
	}

	cs, err := containers(args, gc.recursive)
	if err != nil {
		return err
	}

	var fs []*generator.File
	for _, c := range cs {
		g := generator.New(
			c.Container,
			generator.WithFilesRelativeTo(c.FileRef),
			generator.WithScriptFilename(gc.scriptFilename),
			generator.WithRepoNameBase(gc.repoNameBase),
		)

		gfs, err := g.Files()
		if err != nil {
			return err
		}

		fs = append(fs, gfs...)
	}

	for _, f := range fs {
		dest := relativePath(f.Ref)

		prev, err := ioutil.ReadFile(dest)
		if os.IsNotExist(err) {
			// This is okay since we'll create the file anyway.
		} else if err != nil {
			return err
		}

		if !gc.write {
			diff := difflib.UnifiedDiff{
				A:        splitLinesUnlessEmpty(string(prev)),
				B:        splitLinesUnlessEmpty(f.Content),
				FromFile: dest,
				ToFile:   dest + ".new",
				Context:  3,
			}
			if err := difflib.WriteUnifiedDiff(cmd.OutOrStdout(), diff); err != nil {
				io.WriteString(cmd.ErrOrStderr(), formatError(fmt.Sprintf("Unable to generate changes for file %q", dest), err)+"\n")
				continue
			}
		} else if string(prev) != f.Content {
			if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
				io.WriteString(cmd.ErrOrStderr(), formatError(fmt.Sprintf("Unable to generate file %q", dest), err)+"\n")
				continue
			}

			if err := ioutil.WriteFile(dest, []byte(f.Content), f.Mode); err != nil {
				io.WriteString(cmd.ErrOrStderr(), formatError(fmt.Sprintf("Unable to generate file %q", dest), err)+"\n")
				continue
			}

			if err := os.Chmod(dest, f.Mode); err != nil {
				io.WriteString(cmd.ErrOrStderr(), formatError(fmt.Sprintf("Unable to set permissions for file %q", dest), err)+"\n")
				continue
			}

			if gc.verbose {
				io.WriteString(cmd.OutOrStdout(), dest+"\n")
			}
		}
	}

	return nil
}

func NewGenerateCommand() *cobra.Command {
	gc := &generateCommand{}

	cmd := &cobra.Command{
		Use:     "generate [flags] <paths...>",
		Short:   "Generate the Dockerfiles and build scripts for the given files and directories",
		Aliases: []string{"gen"},
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return gc.run(cmd, args)
		},
	}

	cmd.Flags().BoolVarP(&gc.verbose, "verbose", "v", false, "print each file written")
	cmd.Flags().BoolVarP(&gc.write, "write", "w", false, "whether to persist the generated files")
	cmd.Flags().BoolVarP(&gc.recursive, "recursive", "r", false, "whether to recurse directory arguments")
	cmd.Flags().StringVar(&gc.scriptFilename, "script-filename", generator.DefaultScriptFilename, "the file name to use when generating the build script")
	cmd.Flags().StringVarP(&gc.repoNameBase, "repo-name-base", "n", generator.DefaultRepoNameBase, "the base Docker repository name to use")

	return cmd
}
