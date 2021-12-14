package models

import "dripapp/internal/dripapp/models"

//easyjson:json
type ChatClient struct {
	user models.User
	hub  *Hub
	repo ClientRepository
	io   IOMessage
	send chan Message
}

func NewChatClient(user models.User, hub *Hub, repo ClientRepository, io IOMessage) (client *ChatClient) {
	client = &ChatClient{
		user: user,
		hub:  hub,
		repo: repo,
		io:   io,
		send: make(chan Message),
	}

	hub.register <- client
	return
}

func (c *ChatClient) ReadPump() {
	defer func() {
		c.hub.unregister <- c
	}()

	for {
		var message Message
		err := c.io.ReadMessage(&message)
		if err != nil {
			break
		}

		message, err = c.repo.SaveMessage(c.user.ID, message.ToID, message.Text)
		if err != nil {
			break
		}

		c.hub.broadcast <- message
	}
}

func (c *ChatClient) WritePump() {
	for {
		message := <-c.send
		err := c.io.WriteMessage(message)
		if err != nil {
			break
		}
	}
}

type IOMessage interface {
	ReadMessage(message *Message) error
	WriteMessage(message Message) error
}

type ClientRepository interface {
	SaveMessage(userId uint64, toId uint64, text string) (Message, error)
}
