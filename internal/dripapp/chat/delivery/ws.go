package delivery

import (
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/responses"
	"github.com/gorilla/websocket"
	"net/http"
)

type MessagesWS struct {
	conn   *websocket.Conn
	logger logger.Logger
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *ChatHandler) UpgradeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		status := models.HTTPError{
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
		status := models.HTTPError{
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
