package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func NewLogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "log",
		Short:                 "Submit annotated log messages",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewLogLevelCommand(LogLevelInfo, "Logs a general message"))
	cmd.AddCommand(NewLogLevelCommand(LogLevelWarning, "Logs a warning message"))
	cmd.AddCommand(NewLogLevelCommand(LogLevelError, "Logs an error message"))
	cmd.AddCommand(NewLogLevelCommand(LogLevelFatal, "Logs an error message and forces termination of step container"))

	return cmd
}

type LogLevel string

const (
	LogLevelInfo    LogLevel = "info"
	LogLevelWarning LogLevel = "warn"
	LogLevelError   LogLevel = "error"
	LogLevelFatal   LogLevel = "fatal"
)

func NewLogLevelCommand(level LogLevel, description string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   string(level),
		Short: description,
		Run: func(cmd *cobra.Command, args []string) {
			// This is a stub for eventual forwarding of these log messages to nebula for collection and reporting

			var writer io.Writer

			if level == LogLevelInfo {
				writer = cmd.OutOrStdout()
			} else {
				writer = cmd.ErrOrStderr()
			}

			fmt.Fprintln(
				writer,
				args[0],
			)

			if level == LogLevelFatal {
				os.Exit(1)
			}
		},
	}

	return cmd
}
