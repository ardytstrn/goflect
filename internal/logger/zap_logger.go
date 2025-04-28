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

func (zl *ZapLogger) Warn(msg string, fields ...zap.Field) {
	zl.logger.Warn(msg, fields...)
}

func (zl *ZapLogger) Debug(msg string, fields ...zap.Field) {
	zl.logger.Debug(msg, fields...)
}

func (zl *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	zl.logger.Fatal(msg, fields...)
}

func (zl *ZapLogger) Sync() {
	zl.logger.Sync()
}
