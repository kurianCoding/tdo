package logger

import logr "github.com/sirupsen/logrus"

// Logger common logger for app
type Logger struct {
	l *logr.Logger
}

// Applogger is the instance of Logger
var Applogger Logger

func InitLogger(log *logr.Logger) {
	Applogger.l = log
	return
}

func New() *Logger {
	return &Applogger
}

func (lo *Logger) Info(i ...interface{}) {
	lo.l.Info(i...)
	return
}

func (lo *Logger) Debug(i ...interface{}) {
	lo.l.Debug(i...)
	return
}

func (lo *Logger) Error(i ...interface{}) {
	lo.l.Error(i...)
	return
}
