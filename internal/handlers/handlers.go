package handlers

import (
	"github.com/ardytstrn/goflect/internal/config"
	"github.com/ardytstrn/goflect/internal/logger"
)

type App struct {
	Config config.Config
	Logger logger.Logger
}
