package log

import (
	"fmt"
	"io"
	"os"
)

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// Log prints tagged log messages. It is a stub for future reporting of
// such messages to the nebula service
func writeLog(writer io.Writer, level LogLevel, log string) {
	fmt.Fprintln(
		writer,
		log,
	)
}

// Info reports an informational log
func Info(log string) {
	writeLog(os.Stdout, LogLevelInfo, log)
}

// Warn reports a warning
func Warn(log string) {
	writeLog(os.Stderr, LogLevelWarn, log)
}

// Error reports an error
func Error(log string) {
	writeLog(os.Stderr, LogLevelError, log)
}

// Fatal reports a fatal error then exits the process
func Fatal(log string) {
	writeLog(os.Stderr, LogLevelFatal, log)
	os.Exit(1)
}
