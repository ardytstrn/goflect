package handlers

import (
	"net/url"

	"github.com/ardytstrn/goflect/internal/config"
	"github.com/ardytstrn/goflect/internal/logger"
	"github.com/ardytstrn/goflect/pkg/idgenerator"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Config    config.Config
	Logger    logger.Logger
	Snowflake *idgenerator.Snowflake
	PgPool    *pgxpool.Pool
	Redis     *redis.Client
}

type ErrResponse struct {
	Error string `json:"error"`
}

func isValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}
