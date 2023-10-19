package logger

import (
	"github.com/sirupsen/logrus"
)

type LogrusAdapter struct {
	LogrusLogger *logrus.Logger
}

func NewLogrusAdapter(logrusLogger *logrus.Logger) *LogrusAdapter {
	return &LogrusAdapter{LogrusLogger: logrusLogger}
}

func (a LogrusAdapter) Log(lvl string, msg string, data map[string]interface{}) {
	switch lvl {
	case "debug":
		a.LogrusLogger.WithFields(data).Debug(msg)
	case "info":
		a.LogrusLogger.WithFields(data).Info(msg)
	case "warn":
		a.LogrusLogger.WithFields(data).Warn(msg)
	case "error":
		a.LogrusLogger.WithFields(data).Error(msg)
	case "fatal":
		a.LogrusLogger.WithFields(data).Fatal(msg)
	case "panic":
		a.LogrusLogger.WithFields(data).Panic(msg)
	default:
		a.LogrusLogger.WithFields(data).Error(msg)
	}
}

func (a LogrusAdapter) Debug(msg string, data map[string]interface{}) {
	a.Log("debug", msg, data)
}

func (a LogrusAdapter) Info(msg string, data map[string]interface{}) {
	a.Log("info", msg, data)
}

func (a LogrusAdapter) Warn(msg string, data map[string]interface{}) {
	a.Log("warn", msg, data)
}

func (a LogrusAdapter) Error(msg string, data map[string]interface{}) {
	a.Log("error", msg, data)
}

func (a LogrusAdapter) Panic(msg string, data map[string]interface{}) {
	a.Log("panic", msg, data)
}

func (a LogrusAdapter) Fatal(msg string, data map[string]interface{}) {
	a.Log("fatal", msg, data)
}
