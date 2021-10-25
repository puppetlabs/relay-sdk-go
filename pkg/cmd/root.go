package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() (*cobra.Command, error) {
	c := &cobra.Command{
		Use:           "ni",
		Short:         "Nebula Interface",
		SilenceUsage:  true,
		SilenceErrors: true,
		Long: `The ni tool is meant to be run inside a Relay (FKA "Nebula")
step container to provide helpful SDK-like utilities from the shell.
Invoke it from your relay steps to access parameters and secrets,
generate logs that will be stored on the service, and pass data
between steps.
`,
	}

	c.AddCommand(NewClusterCommand())
	c.AddCommand(NewCredentialsCommand())
	c.AddCommand(NewDecoratorCommand())
	c.AddCommand(NewFileCommand())
	c.AddCommand(NewGetCommand())
	c.AddCommand(NewGitCommand())
	c.AddCommand(NewAWSCommand())
	c.AddCommand(NewOutputCommand())
	c.AddCommand(NewLogCommand())
	c.AddCommand(NewDocCommand())
	c.AddCommand(NewGCPCommand())
	c.AddCommand(NewMetadataCommand())
	c.AddCommand(NewAzureCommand())
	c.AddCommand(NewWorkflowsCommand())

	return c, nil
}
