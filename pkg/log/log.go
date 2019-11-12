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
func writeLog(writer io.Writer, level LogLevel, log ...interface{}) {
	fmt.Fprintln(
		writer,
		log...,
	)
}

// Info reports an informational log
func Info(log string) {
	writeLog(os.Stdout, LogLevelInfo, log)
}

// InfoE reports an informational log from a go error
func InfoE(err error) {
	writeLog(os.Stderr, LogLevelInfo, err)
}

// Warn reports a warning
func Warn(log string) {
	writeLog(os.Stderr, LogLevelWarn, log)
}

// WarnE reports a warning from a go error
func WarnE(err error) {
	writeLog(os.Stderr, LogLevelWarn, err)
}

// Error reports an error
func Error(log string) {
	writeLog(os.Stderr, LogLevelError, log)
}

// ErrorE reports an error from a go error
func ErrorE(err error) {
	writeLog(os.Stderr, LogLevelError, err)
}

// Fatal reports a fatal error then exits the process
func Fatal(log string) {
	writeLog(os.Stderr, LogLevelFatal, log)
	os.Exit(1)
}

// FatalE reports a fatal error from a go error then exits the process
func FatalE(err error) {
	writeLog(os.Stderr, LogLevelFatal, err)
	os.Exit(1)
}
