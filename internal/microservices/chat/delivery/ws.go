package delivery

import (
	_userModels "dripapp/internal/dripapp/models"
	"dripapp/internal/microservices/chat/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/responses"
	"net/http"

	"github.com/gorilla/websocket"
)

type MessagesWS struct {
	conn   *websocket.Conn
	logger logger.Logger
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *ChatHandler) Notifications(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		status := _userModels.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err,
		}
		responses.SendError(w, status, h.Logger.ErrorLogging)
		return
	}

	io := &MessagesWS{
		conn:   conn,
		logger: h.Logger,
	}

	err = h.Chat.ClientHandler(r.Context(), io)
	if err != nil {
		status := _userModels.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}
		responses.SendError(w, status, h.Logger.ErrorLogging)
		return
	}
}

func (m *MessagesWS) ReadMessage(message *models.Message) error {
	err := m.conn.ReadJSON(message)
	if err != nil {
		m.logger.ErrorLogging(http.StatusInternalServerError, "ReadJSON: "+err.Error())
	}

	return err
}

func (m *MessagesWS) WriteMessage(message models.Message) error {
	err := m.conn.WriteJSON(message)
	if err != nil {
		m.logger.ErrorLogging(http.StatusInternalServerError, "WriteJSON: "+err.Error())
	}

	return err
}
