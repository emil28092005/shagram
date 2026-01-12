package websocket

type Hub struct {
	rooms map[string]*Room
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
	}
}

func (h *Hub) GetOrCreateRoom(roomID string) *Room {
	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = NewRoom(roomID)
	}
	return h.rooms[roomID]
}

func (h *Hub) CleanupRoom(roomID string) {
	room := h.rooms[roomID]
	if len(room.clients) == 0 {
		delete(h.rooms, roomID)
	}
}

func (h *Hub) Rooms() []string {
	ids := make([]string, 0, len(h.rooms))
	for id := range h.rooms {
		ids = append(ids, id)
	}
	return ids
}
