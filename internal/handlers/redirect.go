package handlers

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

func (a *App) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")

	var url string

	err := a.PgPool.QueryRow(context.Background(), "SELECT original_url FROM urls WHERE short_code = $1", shortCode).Scan(&url)

	a.Logger.Info("Redirecting!", zap.String("url", url), zap.String("short_code", shortCode))
	if err != nil || shortCode == "" {
		http.Error(w, "404 - The short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
