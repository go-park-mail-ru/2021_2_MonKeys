package delivery

import (
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/responses"
	"net/http"

	"github.com/gorilla/websocket"
)

type WS interface {
	ReadJSON(interface{}) error
	WriteJSON(interface{}) error
}

type NotificationsWS struct {
	conn   WS
	logger logger.Logger
}

func (m *NotificationsWS) Send(user models.User) error {
	err := m.conn.WriteJSON(user)
	if err != nil {
		m.logger.ErrorLogging(http.StatusInternalServerError, "WriteJSON: "+err.Error())
	}

	return err
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *UserHandler) UpgradeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		status := models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err,
		}
		responses.SendError(w, status, h.Logger.ErrorLogging)
		return
	}

	notifications := &NotificationsWS{
		conn:   conn,
		logger: h.Logger,
	}
	err = h.UserUCase.ClientHandler(r.Context(), notifications)
	if err != nil {
		status := models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}
		responses.SendError(w, status, h.Logger.ErrorLogging)
		return
	}
}
