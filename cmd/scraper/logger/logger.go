package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var logger *CompositeLogger

func init() {
	logger = NewLogger(os.Stdout, log.Lmsgprefix|log.Ldate|log.Ltime)
}

type CompositeLogger struct {
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
}

func (l *CompositeLogger) Warning(message string, values ...interface{}) {
	l.WarningLogger.Printf(message, values...)
}

func (l *CompositeLogger) Info(message string, values ...interface{}) {
	l.InfoLogger.Printf(message, values...)
}

func (l *CompositeLogger) Error(message string, values ...interface{}) {
	l.ErrorLogger.Fatalf(fmt.Sprintf(message, values...))
}

func NewLogger(out io.Writer, flags int) *CompositeLogger {
	return &CompositeLogger{
		WarningLogger: log.New(out, "\033[33m[WRN]\033[0m ", flags),
		InfoLogger:    log.New(out, "\033[36m[INF]\033[0m ", flags),
		ErrorLogger:   log.New(out, "\033[31m[ERR]\033[0m ", flags),
	}
}

func Warning(message string, values ...interface{}) {
	logger.Warning(message, values...)
}

func Info(message string, values ...interface{}) {
	logger.Info(message, values...)
}

func Error(message string, values ...interface{}) {
	logger.Error(message, values...)
}
