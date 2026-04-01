package realtime

import (
	"encoding/json"
	"sync"

	"apiRest/internal/entity"
	"github.com/gorilla/websocket"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{clients: map[string]map[*websocket.Conn]struct{}{}}
}

func (h *Hub) Add(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[userID]; !ok {
		h.clients[userID] = map[*websocket.Conn]struct{}{}
	}
	h.clients[userID][conn] = struct{}{}
}

func (h *Hub) Remove(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[userID]; !ok {
		return
	}
	delete(h.clients[userID], conn)
	if len(h.clients[userID]) == 0 {
		delete(h.clients, userID)
	}
}

func (h *Hub) Broadcast(userID string, n entity.Notification) error {
	payload, err := json.Marshal(n)
	if err != nil {
		return err
	}

	h.mu.RLock()
	conns := h.clients[userID]
	h.mu.RUnlock()

	for conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			h.Remove(userID, conn)
			_ = conn.Close()
		}
	}

	return nil
}
