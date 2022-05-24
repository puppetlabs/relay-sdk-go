package cmd

import (
	"time"

	"github.com/puppetlabs/relay-sdk-go/pkg/task"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"github.com/spf13/cobra"
)

const (
	DefaultMetadataTimeout = 5 * time.Minute
)

func NewMetadataCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "metadata",
		Aliases: []string{"meta", "md"},
		Short:   "Manage Relay metadata",
	}

	cmd.AddCommand(NewMetadataRetrieveCommand())

	return cmd
}

func NewMetadataRetrieveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "retrieve",
		Aliases: []string{"get"},
		Short:   "Retrieve Relay metadata",
	}

	cmd.AddCommand(NewMetadataRetrieveEnvironmentCommand())

	return cmd
}

func NewMetadataRetrieveEnvironmentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "environment",
		Aliases: []string{"env"},
		Short:   "Retrieve environment variables",
	}

	cmd.AddCommand(NewMetadataRetrieveEnvironmentVariablesCommand())
	cmd.AddCommand(NewMetadataRetrieveEnvironmentVariableCommand())

	return cmd
}

func NewMetadataRetrieveEnvironmentVariablesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "variables",
		Aliases: []string{"vars"},
		Short:   "Retrieve all environment variables",
		RunE:    doRetrieveEnvironmentVariables,
	}

	return cmd
}

func NewMetadataRetrieveEnvironmentVariableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "variable [name]",
		Aliases: []string{"var"},
		Short:   "Retrieve an environment variable by name",
		Args:    cobra.MaximumNArgs(1),
		RunE:    doRetrieveEnvironmentVariable,
	}

	return cmd
}

func doRetrieveEnvironmentVariables(cmd *cobra.Command, args []string) error {
	planOpts := taskutil.DefaultPlanOptions{}
	t := task.NewTaskInterface(planOpts)

	data, err := t.ReadEnvironmentVariables()
	if err != nil {
		return err
	}

	if _, err = cmd.OutOrStdout().Write(data); err != nil {
		return err
	}

	return nil
}

func doRetrieveEnvironmentVariable(cmd *cobra.Command, args []string) error {
	planOpts := taskutil.DefaultPlanOptions{}
	t := task.NewTaskInterface(planOpts)

	data, err := t.ReadEnvironmentVariable(args[0])
	if err != nil {
		return err
	}

	if _, err = cmd.OutOrStdout().Write(data); err != nil {
		return err
	}

	return nil
}
