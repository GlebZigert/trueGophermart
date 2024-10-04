package logger

import (
	"context"

	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
	ctx    context.Context
}

func NewZapLogger(loggerType string, ctx context.Context) *ZapLogger {
	logger, _ := zap.NewProduction()

	return &ZapLogger{logger: logger, ctx: ctx}
}

func (l *ZapLogger) Debug(msg string, fields map[string]interface{}) {

	l.logger.Debug("", zap.Any("args", fields))
}

func (l *ZapLogger) Info(msg string, fields map[string]interface{}) {

	l.logger.Info(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Warn(msg string, fields map[string]interface{}) {

	l.logger.Warn("", zap.Any("args", fields))
}

func (l *ZapLogger) Error(msg string, fields map[string]interface{}) {

	l.logger.Error(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Fatal(msg string, fields map[string]interface{}) {

	l.logger.Fatal("", zap.Any("args", fields))
}
