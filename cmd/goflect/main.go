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
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	godotenv.Load()

	cfg := config.Load()

	logger, _ := setupLogger(cfg.Environment)

	logger.Debug("Config",
		zap.String("environment", cfg.Environment),
		zap.String("postgresURL", cfg.PostgresURL),
		zap.String("redisURL", cfg.RedisURL),
		zap.String("port", cfg.Port),
	)

	logger.Info("Started in " + cfg.Environment + " environment")

	defer logger.Sync()

	// Initialize PostgreSQL connection pool
	pgConfig, err := pgxpool.ParseConfig(cfg.PostgresURL)
	if err != nil {
		logger.Fatal("Failed to parse Postgres URL", zap.Error(err))
		return
	}
	pgConfig.MaxConns = 100

	pgPool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		logger.Fatal("Failed to parse Postgres connection config", zap.Error(err))
		return
	}

	_, err = pgPool.Acquire(context.Background())
	if err != nil {
		logger.Fatal("Postgres connection failed", zap.Error(err))
		return
	}

	logger.Info("Connected to PostgreSQL")

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
		DB:   0,
	})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatal("Redis connection failed", zap.Error(err))
		return
	}

	logger.Info("Connected to Redis")

	app := &handlers.App{
		Config:    cfg,
		Logger:    logger,
		Snowflake: idgenerator.NewSnowflake(1),
		PgPool:    pgPool,
		Redis:     rdb,
	}

	app.StartBatchWorkers()

	router := http.NewServeMux()
	router.HandleFunc("/api/shorten", app.ShortenHandler)
	router.HandleFunc("/{shortCode}", app.RedirectHandler)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.Chain(router, app),
	}

	logger.Info("Starting server", zap.String("port", cfg.Port))

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Server error", zap.Error(err))
	}
}

func setupLogger(environment string) (*logger.ZapLogger, error) {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.DisableCaller = true
	config.EncoderConfig.EncodeCaller = nil

	if environment != "production" {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.DisableStacktrace = false
		config.Sampling = nil
	}

	zl, err := config.Build()

	if err != nil {
		return nil, err
	}

	return logger.NewZapLogger(zl), nil
}
