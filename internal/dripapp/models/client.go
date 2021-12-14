package models

type Client struct {
	user          User
	hub           *Hub
	notifications Notifications
}

func NewClient(user User, hub *Hub, notifications Notifications) (client *Client) {
	client = &Client{
		user:          user,
		hub:           hub,
		notifications: notifications,
	}

	hub.register <- client
	return
}

func (c *Client) NotifyAboutMatchWith(user User) {
	err := c.notifications.Send(user)
	if err != nil {
		c.hub.unregister <- c
	}
}
