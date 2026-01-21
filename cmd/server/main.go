package main

import (
	"log"
	"net/http"
	"os"
	"shagram/internal/api"
	"shagram/internal/auth"
	"shagram/internal/db"
	"shagram/internal/models"
	"shagram/internal/websocket"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "/app/data/shagram.db"
	}
	database, err := db.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	hub := websocket.NewHub()
	hub.GetOrCreateRoom("general")
	hub.GetOrCreateRoom("chat")
	hub.GetOrCreateRoom("dev")

	router := gin.Default()

	router.POST("/api/auth/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
			return
		}

		token, err := auth.NewAccessToken(req.Username, time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"access_token": token})
	})
	router.GET("/api/me", auth.Middleware(), func(c *gin.Context) {
		username := c.GetString(auth.CtxUsernameKey)
		c.JSON(200, gin.H{"username": username})
	})
	router.GET("/ws/:room", func(c *gin.Context) {
		api.WebSocketHandler(hub, database)(c)
	})
	router.GET("/api/rooms", func(c *gin.Context) {
		rows, err := database.Query(`
			SELECT DISTINCT room_id
			FROM messages
			ORDER BY room_id`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var rooms []string
		for rows.Next() {
			var roomID string
			err := rows.Scan(&roomID)
			if err != nil {
				continue
			}
			rooms = append(rooms, roomID)
		}

		if len(rooms) == 0 {
			rooms = []string{"general", "chat", "dev"}
		}

		c.JSON(200, gin.H{"rooms": rooms})
	})
	router.GET("/api/messages/:room", func(c *gin.Context) {
		roomID := c.Param("room")
		rows, err := database.Query(`
			SELECT id, room_id, user, text, created_at
			FROM messages
			WHERE room_id = ?
			ORDER BY created_at DESC
			LIMIT 50`, roomID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		var messages []models.Message
		for rows.Next() {
			var msg models.Message
			err := rows.Scan(&msg.ID, &msg.RoomID, &msg.User, &msg.Text, &msg.CreatedAt)
			if err != nil {
				log.Printf("Scan error: %v", err)
				continue
			}
			messages = append(messages, msg)
		}

		c.JSON(200, gin.H{"messages": messages})
	})
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	router.Run(":8080")
}
