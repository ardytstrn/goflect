package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ardytstrn/goflect/pkg/util"
	"github.com/bytedance/sonic"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Task struct {
	ShortCode   string
	OriginalURL string
}

var taskQueue = make(chan Task, 10000)

func (a *App) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	dec := sonic.ConfigFastest.NewDecoder(r.Body)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req struct {
		URL string `json:"url"`
	}

	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResp(w, ErrResponse{Error: "Invalid JSON"})
		return
	}

	if !isValidURL(req.URL) {
		w.WriteHeader(http.StatusBadRequest)
		errorResp(w, ErrResponse{Error: "Invalid URL"})
		return
	}

	id := a.Snowflake.Generate()
	shortCode := util.EncodeBase62(id)

	taskQueue <- Task{ShortCode: shortCode, OriginalURL: req.URL}

	response := struct {
		ShortURL string `json:"short_url"`
	}{
		ShortURL: fmt.Sprintf("https://%s/%s", a.Config.Domain, shortCode),
	}

	respBytes, err := sonic.Marshal(&response)
	if err != nil {
		a.Logger.Error("Failed to marshal response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		errorResp(w, ErrResponse{Error: "Internal server error"})
		return
	}

	w.Write(respBytes)
}

func (a *App) StartBatchWorkers() {
	const batchSize = 100
	const flushInterval = 100 * time.Millisecond

	for i := 0; i < 10; i++ {
		go func() {
			var batch []Task
			timer := time.NewTimer(flushInterval)

			for {
				select {
				case task := <-taskQueue:
					batch = append(batch, task)
					if len(batch) >= batchSize {
						a.flushBatch(batch)
						batch = nil
						if !timer.Stop() {
							<-timer.C
						}
						timer.Reset(flushInterval)
					}
				case <-timer.C:
					if len(batch) > 0 {
						a.flushBatch(batch)
						batch = nil
					}
					timer.Reset(flushInterval)
				}
			}
		}()
	}
}

func (a *App) flushBatch(tasks []Task) {
	ctx := context.Background()
	inputRows := make([][]interface{}, len(tasks))

	for i, task := range tasks {
		inputRows[i] = []interface{}{task.ShortCode, task.OriginalURL}
	}

	_, err := a.PgPool.CopyFrom(
		ctx,
		pgx.Identifier{"urls"},
		[]string{"short_code", "original_url"},
		pgx.CopyFromRows(inputRows),
	)

	if err != nil {
		a.Logger.Error("COPY failed", zap.Error(err))
	}
}

func errorResp(w http.ResponseWriter, resp ErrResponse) {
	respBytes, _ := sonic.Marshal(&resp)
	w.Write(respBytes)
}
