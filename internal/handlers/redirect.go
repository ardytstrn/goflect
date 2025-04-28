package handlers

import (
	"net/http"

	"go.uber.org/zap"
)

func (a *App) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortCode := r.PathValue("shortCode")

	if shortCode == "" {
		http.Error(w, "Invalid short code", http.StatusBadRequest)
		return
	}

	url, err := a.Redis.Get(ctx, shortCode).Result()

	if err == nil {
		a.Logger.Debug("Cache hit for short code", zap.String("code", shortCode))
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	a.Logger.Info("Cache miss or Redis error", zap.String("code", shortCode), zap.Error(err))

	err = a.PgPool.QueryRow(ctx, "SELECT original_url FROM urls WHERE short_code = $1", shortCode).Scan(&url)
	if err != nil {
		http.Error(w, "404 - The short URL not found", http.StatusNotFound)
		return
	}

	if err := a.Redis.Set(ctx, shortCode, url, 0).Err(); err != nil {
		a.Logger.Warn("Failed to cache short code after DB lookup",
			zap.String("code", shortCode),
			zap.Error(err))
	} else {
		a.Logger.Debug("Cached short code after DB lookup", zap.String("code", shortCode))
	}

	a.Logger.Info("Redirecting to original URL", zap.String("url", url), zap.String("code", shortCode))
	http.Redirect(w, r, url, http.StatusFound)
}
