package http

import (
	"log/slog"
	"net/http"
	"time"

	"apiRest/internal/infra/realtime"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	Hub    *realtime.Hub
	Logger *slog.Logger
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WebSocketHandler) Handle(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user in token"})
		return
	}

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		if h.Logger != nil {
			h.Logger.Error("websocket upgrade failed", "error", err)
		}
		return
	}

	h.Hub.Add(userID, conn)
	defer func() {
		h.Hub.Remove(userID, conn)
		_ = conn.Close()
	}()

	for {
		_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}
