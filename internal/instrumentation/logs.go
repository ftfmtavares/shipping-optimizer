package instrumentation

import (
	"log"
	"os"
)

type Logger struct {
	infoLog    *log.Logger
	warningLog *log.Logger
	errorLog   *log.Logger
}

func NewLogger() Logger {
	return Logger{
		infoLog:    log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime),
		warningLog: log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Ltime),
		errorLog:   log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime),
	}
}

func (l Logger) Info(msg string) {
	l.infoLog.Println(" " + msg)
}

func (l Logger) Warning(msg string) {
	l.warningLog.Println(" " + msg)
}

func (l Logger) Error(msg string) {
	l.errorLog.Println(" " + msg)
}
