package models

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			_, ok := h.clients[client]
			if ok {
				delete(h.clients, client)
			}
		}
	}
}

func (h *Hub) NotifyAboutMatchWith(recipientId uint64, user User) {
	for client := range h.clients {
		if client.user.ID == recipientId {
			err := client.notifications.Send(user)
			if err != nil {
				h.unregister <- client
			}
		}
	}
}
