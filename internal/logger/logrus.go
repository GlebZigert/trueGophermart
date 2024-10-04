package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	logger *logrus.Logger
	ctx    context.Context
}

func NewLogrusLogger(loggerType string, ctx context.Context) *LogrusLogger {
	logger := logrus.New()
	logger.Out = os.Stdout

	return &LogrusLogger{logger: logger, ctx: ctx}
}

func (l *LogrusLogger) Debug(msg string, fields map[string]interface{}) {

	l.logger.WithFields(fields).Debug(msg)
}

func (l *LogrusLogger) Info(msg string, fields map[string]interface{}) {

	l.logger.WithFields(fields).Info(msg)
}

func (l *LogrusLogger) Warn(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Warn(msg)
}

func (l *LogrusLogger) Error(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Error(msg)
}

func (l *LogrusLogger) Fatal(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Fatal(msg)
}
