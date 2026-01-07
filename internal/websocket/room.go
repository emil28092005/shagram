package websocket

import "github.com/gorilla/websocket"

type Client struct {
	Conn *websocket.Conn
	Room *Room
}

type Room struct {
	id      string
	clients map[*Client]bool
}

func NewRoom(id string) *Room {
	return &Room{
		id:      id,
		clients: make(map[*Client]bool),
	}
}

func (r *Room) Register(client *Client) {
	r.clients[client] = true
}

func (r *Room) Unregister(client *Client) {
	delete(r.clients, client)
	client.Conn.Close()
}

func (r *Room) Broadcast(message []byte) {
	for client := range r.clients {
		err := client.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			r.Unregister(client)
		}
	}
}
