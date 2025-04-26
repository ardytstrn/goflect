package logger

import "go.uber.org/zap"

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(zapLogger *zap.Logger) *ZapLogger {
	return &ZapLogger{logger: zapLogger}
}

func (zl *ZapLogger) Info(msg string, fields ...zap.Field) {
	zl.logger.Info(msg, fields...)
}

func (zl *ZapLogger) Error(msg string, fields ...zap.Field) {
	zl.logger.Error(msg, fields...)
}

func (zl *ZapLogger) Sync() {
	zl.logger.Sync()
}
