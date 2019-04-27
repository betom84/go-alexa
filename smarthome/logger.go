package smarthome

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Logger interface
type Logger interface {
	Trace(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

// Severe level to create new default logger
const (
	Trace = iota
	Debug
	Info
	Warning
	Error
	Fatal
)

type defaultLogger struct {
	Level  int
	stdLog *log.Logger
}

// Log is the Logger implementation got used for logging
var Log = NewDefaultLogger(Error, os.Stderr)

// NewDefaultLogger creates a default Logger implemtation uses stdlib log with default flags
func NewDefaultLogger(level int, out io.Writer) Logger {
	return &defaultLogger{
		Level:  level,
		stdLog: log.New(out, "", log.LstdFlags),
	}
}

// Trace message
func (logger defaultLogger) Trace(format string, v ...interface{}) {
	if logger.Level > Trace {
		return
	}

	logger.log("TRACE", format, v...)
}

// Debug message
func (logger defaultLogger) Debug(format string, v ...interface{}) {
	if logger.Level > Debug {
		return
	}

	logger.log("DEBUG", format, v...)
}

// Info message
func (logger defaultLogger) Info(format string, v ...interface{}) {
	if logger.Level > Info {
		return
	}

	logger.log("INFO", format, v...)
}

// Warning message
func (logger defaultLogger) Warning(format string, v ...interface{}) {
	if logger.Level > Warning {
		return
	}

	logger.log("WARN", format, v...)
}

// Error message
func (logger defaultLogger) Error(format string, v ...interface{}) {
	if logger.Level > Error {
		return
	}

	logger.log("ERROR", format, v...)
}

// Fatal message
func (logger defaultLogger) Fatal(format string, v ...interface{}) {
	logger.log("FATAL", format, v...)
}

func (logger defaultLogger) log(level, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	if strings.HasSuffix(message, "\n") {
		message = message[:len(message)-1]
	}

	logger.stdLog.Printf("[go-alexa] [%s] %s\n", level, message)
}
