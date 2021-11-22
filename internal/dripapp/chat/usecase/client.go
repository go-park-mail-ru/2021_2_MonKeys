package usecase

import (
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"github.com/gorilla/websocket"
	"net/http"
)

type Client struct {
	user models.User
	uc   *ChatUseCase
	conn *websocket.Conn
	send chan models.Message
}

func (c *Client) readPump() {
	defer func() {
		c.uc.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var message models.Message
		err := c.conn.ReadJSON(&message)
		if err != nil {
			logger.DripLogger.ErrorLogging(http.StatusInternalServerError, "ReadJSON: "+err.Error())
			break
		}

		message, err = c.uc.SaveMessage(c.user, message)
		if err != nil {
			logger.DripLogger.ErrorLogging(http.StatusInternalServerError, "UseCase: "+err.Error())
			break
		}

		c.uc.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		message := <-c.send
		err := c.conn.WriteJSON(message)
		if err != nil {
			logger.DripLogger.ErrorLogging(http.StatusInternalServerError, "WriteJSON: "+err.Error())
			return
		}
	}
}
