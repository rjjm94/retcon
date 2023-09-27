package logger

import (
	"fmt"
	"log"
	"os"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	*log.Logger
	level Level
	file  *os.File
}

func NewLogger(logFileName string, level Level) (*Logger, error) {
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	l := log.New(logFile, "", log.LstdFlags)
	return &Logger{l, level, logFile}, nil
}

func (l *Logger) Close() error {
	if err := l.file.Close(); err != nil {
		return fmt.Errorf("failed to close log file: %v", err)
	}
	return nil
}

func (l *Logger) Debug(msg string) {
	if l.level <= DEBUG {
		l.Output(2, fmt.Sprintf("DEBUG: %s", msg))
	}
}

func (l *Logger) Info(msg string) {
	if l.level <= INFO {
		l.Output(2, fmt.Sprintf("INFO: %s", msg))
	}
}

func (l *Logger) Warn(msg string) {
	if l.level <= WARN {
		l.Output(2, fmt.Sprintf("WARN: %s", msg))
	}
}

func (l *Logger) Error(msg string) {
	if l.level <= ERROR {
		l.Output(2, fmt.Sprintf("ERROR: %s", msg))
	}
}

func (l *Logger) Fatal(msg string) {
	if l.level <= FATAL {
		l.Output(2, fmt.Sprintf("FATAL: %s", msg))
		os.Exit(1)
	}
}
