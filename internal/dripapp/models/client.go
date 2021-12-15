package models

import "fmt"

type Client struct {
	user          User
	hub           *Hub
	notifications Notifications
}

func NewClient(user User, hub *Hub, notifications Notifications) {
	fmt.Println("---------6----------")
	client := &Client{
		user:          user,
		hub:           hub,
		notifications: notifications,
	}
	fmt.Println("---------7----------")

	hub.register <- client
	fmt.Println("---------8----------")
	return
}

func (c *Client) NotifyAboutMatchWith(user User) {
	fmt.Println(user)
	err := c.notifications.Send(user)
	fmt.Println("---------89----------")
	if err != nil {
		fmt.Println("---------80----------")
		c.hub.unregister <- c
	}
	fmt.Println("---------8q----------")
}
