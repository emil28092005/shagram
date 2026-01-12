package api

import (
	"log"
	"net/http"
	"shagram/internal/db"
	"shagram/internal/websocket"

	"github.com/gin-gonic/gin"
	gws "github.com/gorilla/websocket"
)

var upgrader = gws.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(hub *websocket.Hub, database *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("room")
		room := hub.GetOrCreateRoom(roomID)

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
				err := client.Conn.ReadJSON(&msg)
				if err != nil {
					break
				}
				_, err = database.Exec(`
					INSERT INTO messages (room_id, user, text)
					VALUES (?, ?, ?)`, roomID, "user", msg["text"])
				if err != nil {
					log.Printf("Save message error: %v", err)
				}

				message := []byte(msg["text"])
				room.Broadcast(message)
			}
		}()
	}

}
