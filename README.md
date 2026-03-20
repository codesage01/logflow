# LogFlow 🚀
> Real-time Log Aggregation & Monitoring System — Built in Go

A production-style log aggregation service with REST API, WebSocket streaming, and a live dashboard. Built to demonstrate real-world Go patterns.

## Features
- **REST API** — ingest, query, and aggregate logs
- **Real-time streaming** — WebSocket broadcasts every new log instantly
- **Concurrent-safe storage** — `sync.RWMutex` for thread-safe read/write
- **Live dashboard** — browser UI with filters, search, and per-level stats
- **Interface-driven storage** — swap MemoryStore → PostgreSQL with zero handler changes

## Project Structure
```
logflow/
├── cmd/server/main.go          # Entry point, router setup
├── config/config.go            # Env-based config
├── internal/
│   ├── models/log.go           # LogEntry, QueryFilter, Stats types
│   ├── storage/memory.go       # Thread-safe in-memory store (Store interface)
│   ├── hub/hub.go              # WebSocket broadcast hub (goroutine-based)
│   └── handlers/
│       ├── log_handler.go      # POST /api/logs, GET /api/logs, GET /api/logs/stats
│       └── ws_handler.go       # GET /ws — WebSocket upgrade + streaming
└── web/index.html              # Live dashboard frontend

## Quick Start
```bash
go mod tidy
go run ./cmd/server

# Server starts on :8080

```
![LogFlow Dashboard](./assets/dashboard.png)
## API

### Ingest a log
```bash
curl -X POST http://localhost:8080/api/logs \
  -H "Content-Type: application/json" \
  -d '{"level":"ERROR","service":"auth-service","message":"JWT validation failed","meta":{"user_id":"u_123"}}'
```

### Query logs
```bash
curl "http://localhost:8080/api/logs?level=ERROR&service=auth-service&search=JWT&limit=50"
```

### Get stats
```bash
curl http://localhost:8080/api/logs/stats
```

### Live dashboard
Open `http://localhost:8080` in your browser.

## Go Concepts Demonstrated
| Concept | Where |
|---|---|
| Goroutines | `hub.Run()`, WebSocket pump goroutines |
| Channels | `hub.broadcast`, `hub.register`, `client.Send` |
| `sync.RWMutex` | `MemoryStore` concurrent read/write |
| Interfaces | `storage.Store` — swap implementations easily |
| Context & select | Hub event loop with `select` |
| HTTP routing (stdlib) | `http.NewServeMux()` with method+path patterns |
| WebSocket upgrade | `gorilla/websocket` Upgrader |
| JSON encoding | `encoding/json` throughout |



