package logger

import (
	"log"
	"time"
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

func (l Logger) writeMsg(msg string) {
	log.Printf("[%s] %s: %s", time.Now().Format(time.RFC822), l.level, msg)
}

func (l Logger) Info(msg string) {
	if logLevelAllowed(l.level, info) {
		l.writeMsg(msg)
	}
}

func (l Logger) Error(msg string) {
	if logLevelAllowed(l.level, err) {
		l.writeMsg(msg)
	}
}

func (l Logger) Debug(msg string) {
	if logLevelAllowed(l.level, debug) {
		l.writeMsg(msg)
	}
}

func (l Logger) Warn(msg string) {
	if logLevelAllowed(l.level, warn) {
		l.writeMsg(msg)
	}
}
