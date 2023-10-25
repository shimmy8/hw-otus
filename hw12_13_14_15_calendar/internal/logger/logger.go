package logger

import (
	"log"
)

const (
	debug = "DEBUG"
	info  = "INFO"
	warn  = "WARN"
	err   = "ERROR"
)

var levelPriority = map[string]int{
	debug: 0,
	info:  1,
	warn:  2,
	err:   3,
}

func logLevelAllowed(logLevel string, msgLevel string) bool {
	return levelPriority[logLevel] <= levelPriority[msgLevel]
}

type Logger struct {
	level string
	name  string
}

func New(level string, name string) *Logger {
	return &Logger{level: level, name: name}
}

func (l Logger) Info(msg string) {
	if logLevelAllowed(l.level, info) {
		log.Println(msg)
	}
}

func (l Logger) Error(msg string) {
	if logLevelAllowed(l.level, err) {
		log.Println(msg)
	}
}

func (l Logger) Debug(msg string) {
	if logLevelAllowed(l.level, debug) {
		log.Println(msg)
	}
}

func (l Logger) Warn(msg string) {
	if logLevelAllowed(l.level, warn) {
		log.Println(msg)
	}
}
