// Package instrumentation handles logging
package instrumentation

import (
	"log"
	"os"
)

// Logger provides logging methods
type Logger struct {
	infoLog    *log.Logger
	warningLog *log.Logger
	errorLog   *log.Logger
}

// NewLogger initializes a Logger
func NewLogger() Logger {
	return Logger{
		infoLog:    log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime),
		warningLog: log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Ltime),
		errorLog:   log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime),
	}
}

// Info method logs an info event
func (l Logger) Info(msg string) {
	l.infoLog.Println(msg)
}

// Warning method logs a warning event
func (l Logger) Warning(msg string) {
	l.warningLog.Println(msg)
}

// Error method logs an error event
func (l Logger) Error(msg string) {
	l.errorLog.Println(msg)
}
