package models

type Client struct {
	user          User
	hub           *Hub
	notifications Notifications
}

func NewClient(user User, hub *Hub, notifications Notifications) {
	client := &Client{
		user:          user,
		hub:           hub,
		notifications: notifications,
	}

	hub.register <- client
}
