package main

import (
	"net/http"

	"github.com/ardytstrn/goflect/internal/config"
	"github.com/ardytstrn/goflect/internal/handlers"
	"github.com/ardytstrn/goflect/internal/middleware"
)

func main() {
	cfg := config.Load()

	router := http.NewServeMux()

	app := &handlers.App{
		Config: cfg,
	}

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.Chain(router, app),
	}

	server.ListenAndServe()
}
