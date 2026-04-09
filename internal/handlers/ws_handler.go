package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/codesage01/logflow/internal/hub"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for dev — restrict in production
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebSocketHandler struct {
	hub *hub.Hub
}

func NewWebSocketHandler(h *hub.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: h}
}

// Handle upgrades HTTP to WebSocket and streams logs to the client
func (h *WebSocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &hub.Client{Send: make(chan []byte, 64)}
	h.hub.Register(client)

	// Write pump — sends messages from hub to WebSocket client
	go func() {
		defer func() {
			h.hub.Unregister(client)
			conn.Close()
		}()

		for msg := range client.Send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				break
			}
		}
	}()

	// Read pump — keeps connection alive, detects client disconnect
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.hub.Unregister(client)
			break
		}
	}
}
