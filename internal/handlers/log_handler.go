package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/codesage01/logflow/internal/hub"
	"github.com/codesage01/logflow/internal/models"
	"github.com/codesage01/logflow/internal/storage"
)

type LogHandler struct {
	store storage.Store
	hub   *hub.Hub
}

func NewLogHandler(store storage.Store, hub *hub.Hub) *LogHandler {
	return &LogHandler{store: store, hub: hub}
}

// Ingest accepts a new log entry via POST /api/logs
func (h *LogHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	var entry models.LogEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Set server-side fields
	entry.ID = uuid.NewString()
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	if entry.Level == "" {
		entry.Level = models.LevelInfo
	}
	if entry.Service == "" {
		entry.Service = "unknown"
	}

	if err := h.store.Save(&entry); err != nil {
		http.Error(w, "failed to save log", http.StatusInternalServerError)
		return
	}

	// Broadcast to all WebSocket clients in real-time
	h.hub.Broadcast(&entry)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

// Query returns filtered logs via GET /api/logs
func (h *LogHandler) Query(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := models.QueryFilter{
		Level:   models.Level(q.Get("level")),
		Service: q.Get("service"),
		Search:  q.Get("search"),
	}

	logs, err := h.store.Query(filter)
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// Stats returns aggregate statistics via GET /api/logs/stats
func (h *LogHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.Stats()
	if err != nil {
		http.Error(w, "stats failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
