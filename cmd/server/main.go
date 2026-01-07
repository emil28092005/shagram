package main

import (
	"log"
	"shagram/internal/api"
	"shagram/internal/db"
	"shagram/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	database, err := db.NewDB("shagram.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	hub := websocket.NewHub()
	router := gin.Default()

	router.GET("/ws/:room", api.WebSocketHandler(hub))
	router.GET("/api/rooms", func(c *gin.Context) {
		c.JSON(200, gin.H{"rooms": "TODO"})
	})
	router.Run(":8080")
}
