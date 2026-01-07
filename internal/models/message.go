package models

import "time"

type Message struct {
	ID        int       `json:"id"`
	RoomID    string    `json:"room_id"`
	User      string    `json:"user"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
