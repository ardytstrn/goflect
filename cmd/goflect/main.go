package main

import (
	"net/http"

	"github.com/ardytstrn/goflect/internal/config"
	"github.com/ardytstrn/goflect/internal/handlers"
	"github.com/ardytstrn/goflect/internal/logger"
	"github.com/ardytstrn/goflect/internal/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg := config.Load()

	logger, _ := setupLogger()
	defer logger.Sync()

	router := http.NewServeMux()

	app := &handlers.App{
		Config: cfg,
		Logger: logger,
	}

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.Chain(router, app),
	}

	logger.Info("Starting server", zap.String("port", cfg.Port))
	server.ListenAndServe()
}

func setupLogger() (*logger.ZapLogger, error) {
	prodConfig := zap.NewProductionConfig()
	prodConfig.Encoding = "console"
	prodConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	prodConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	prodConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	zl, err := prodConfig.Build()

	if err != nil {
		return nil, err
	}

	return logger.NewZapLogger(zl), nil
}
