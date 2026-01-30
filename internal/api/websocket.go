package api

import (
	"log"
	"net/http"
	"os"
	"shagram/internal/auth"
	"shagram/internal/db"
	"shagram/internal/websocket"
	"strings"

	"github.com/gin-gonic/gin"
	gws "github.com/gorilla/websocket"
)

var upgrader = gws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return false
		}
		allowed := strings.Split(os.Getenv("WS_ALLOWED_ORIGINS"), ",")
		for _, a := range allowed {
			if strings.TrimSpace(a) == origin {
				return true
			}
		}
		return false
	},
}

func WebSocketHandler(hub *websocket.Hub, database *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("room")
		room := hub.GetOrCreateRoom(roomID)

		tokenString := c.Query("token")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := auth.ParseAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		username := claims.Username

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		client := &websocket.Client{Conn: conn, Room: room}
		room.Register(client)

		go func() {
			defer func() {
				room.Unregister(client)
				hub.CleanupRoom(roomID)
			}()

			for {
				var msg map[string]string
				if err := client.Conn.ReadJSON(&msg); err != nil {
					break
				}

				text := msg["text"]
				if text == "" {
					continue
				}

				_, err = database.Exec(`
					INSERT INTO messages (room_id, user, text)
					VALUES (?, ?, ?)`, roomID, username, text)
				if err != nil {
					log.Printf("Save message error: %v", err)
				}

				room.Broadcast([]byte(username + ": " + text))
			}
		}()
	}
}
