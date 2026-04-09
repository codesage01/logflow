
package main

import (
	"log"
	"net/http"

	"github.com/codesage01/logflow/config"
	"github.com/codesage01/logflow/internal/handlers"
	"github.com/codesage01/logflow/internal/hub"
	"github.com/codesage01/logflow/internal/storage"
)

func main() {
	cfg := config.Load()

	// Initialize in-memory storage (swap with PostgreSQL for production)
	store := storage.NewMemoryStore()

	// Initialize WebSocket hub (broadcasts logs to connected clients)
	wsHub := hub.NewHub()
	go wsHub.Run()

	// Setup HTTP router
	mux := http.NewServeMux()

	logHandler := handlers.NewLogHandler(store, wsHub)
	wsHandler := handlers.NewWebSocketHandler(wsHub)

	// REST API routes
	mux.HandleFunc("POST /api/logs", logHandler.Ingest)       // ingest a log
	mux.HandleFunc("GET /api/logs", logHandler.Query)          // query logs
	mux.HandleFunc("GET /api/logs/stats", logHandler.Stats)    // get stats

	// WebSocket route
	mux.HandleFunc("GET /ws", wsHandler.Handle)

	// Serve the frontend dashboard
	mux.Handle("/", http.FileServer(http.Dir("web")))

// Ensure the port has a leading colon
    address := ":" + cfg.Port
    
    log.Printf("LogFlow server running on %s", address)
    
    // Use the address variable which now contains ":8080"
    if err := http.ListenAndServe(address, mux); err != nil {
        log.Fatal(err)
    }
}
