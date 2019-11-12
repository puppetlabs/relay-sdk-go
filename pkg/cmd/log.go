package cmd

import (
	"strings"

	"github.com/puppetlabs/nebula-sdk/pkg/log"
	"github.com/spf13/cobra"
)

func NewLogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "log",
		Short:                 "Submit annotated log messages",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewLogInfoCommand())
	cmd.AddCommand(NewLogWarnCommand())
	cmd.AddCommand(NewLogErrorCommand())
	cmd.AddCommand(NewLogFatalCommand())

	return cmd
}

func NewLogInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   string(log.LogLevelInfo),
		Short: "Logs an informational message",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.Info(strings.Join(args, " "))
		},
	}

	return cmd
}

func NewLogWarnCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   string(log.LogLevelWarn),
		Short: "Logs a warning message",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.Warn(strings.Join(args, " "))
		},
	}

	return cmd
}

func NewLogErrorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   string(log.LogLevelError),
		Short: "Logs an error message",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.Error(strings.Join(args, " "))
		},
	}

	return cmd
}

func NewLogFatalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   string(log.LogLevelFatal),
		Short: "Logs a fatal error message and exits process",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.Fatal(strings.Join(args, " "))
		},
	}

	return cmd
}
