package main

import (
	"context"
	"net/http"

	"github.com/ardytstrn/goflect/internal/config"
	"github.com/ardytstrn/goflect/internal/handlers"
	"github.com/ardytstrn/goflect/internal/logger"
	"github.com/ardytstrn/goflect/internal/middleware"
	"github.com/ardytstrn/goflect/pkg/idgenerator"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg := config.Load()

	logger, _ := setupLogger()

	defer logger.Sync()

	// Initialize PostgreSQL connection pool
	pgConfig, err := pgxpool.ParseConfig(cfg.PostgresURL)
	if err != nil {
		logger.Error("Failed to parse Postgres URL", zap.Error(err))
	}
	pgConfig.MaxConns = 100

	pgPool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		logger.Error("Cannot parse config", zap.Error(err))
		return
	}

	_, err = pgPool.Acquire(context.Background())
	if err != nil {
		logger.Error("Cannot connect PostgreSQL server", zap.Error(err))
		return
	}

	app := &handlers.App{
		Config:    cfg,
		Logger:    logger,
		Snowflake: idgenerator.NewSnowflake(1),
		PgPool:    pgPool,
	}
	app.StartBatchWorkers()

	router := http.NewServeMux()
	router.HandleFunc("/api/shorten", app.ShortenHandler)

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
