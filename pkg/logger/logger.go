package logger

import (
	"fmt"
	"io"
	"os"
)

type Logger interface {
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Success(format string, v ...interface{})
	Log(format string, v ...interface{})
}

type SimpleLogger struct {
	output io.Writer
	errOut io.Writer
}

func NewLogger(output, errOut io.Writer) *SimpleLogger {
	return &SimpleLogger{output: output, errOut: errOut}
}

var defaultLogger Logger = NewLogger(os.Stdout, os.Stderr)

func (l *SimpleLogger) log(writer io.Writer, color, format string, v ...interface{}) {
	coloredFormat := fmt.Sprintf("%s%s%s", color, format, "\033[0m")
	fmt.Fprintf(writer, coloredFormat+"\n", v...)
}

func (l *SimpleLogger) Info(format string, v ...interface{}) {
	l.log(l.output, "\033[34m", format, v...)
}

func (l *SimpleLogger) Warn(format string, v ...interface{}) {
	l.log(l.output, "\033[33m", format, v...)
}

func (l *SimpleLogger) Error(format string, v ...interface{}) {
	l.log(l.errOut, "\033[31m", format, v...)
}

func (l *SimpleLogger) Success(format string, v ...interface{}) {
	l.log(l.output, "\033[32m", format, v...)
}

func (l *SimpleLogger) Log(format string, v ...interface{}) {
	l.log(l.output, "", format, v...)
}

// Package-level functions
func Info(format string, v ...interface{}) {
	defaultLogger.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	defaultLogger.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	defaultLogger.Error(format, v...)
}

func Success(format string, v ...interface{}) {
	defaultLogger.Success(format, v...)
}

func Log(format string, v ...interface{}) {
	defaultLogger.Log(format, v...)
}
