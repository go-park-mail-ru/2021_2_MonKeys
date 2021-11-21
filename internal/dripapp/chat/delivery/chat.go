package delivery

import (
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/responses"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

type ChatHandler struct {
	Chat   models.ChatUseCase
	Logger logger.Logger
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *ChatHandler) Notifications(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		status := models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err,
		}
		responses.SendError(w, status, h.Logger.ErrorLogging)
		return
	}

	currentUser, ok := r.Context().Value(configs.ContextUser).(models.User)
	if !ok {
		status := models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err,
		}
		responses.SendError(w, status, h.Logger.ErrorLogging)
	}

	go h.sendNewMsgNotifications(currentUser, ws)
}

func (h *ChatHandler) sendNewMsgNotifications(currentUser models.User, client *websocket.Conn) {
	for {
		var msg models.Message

		err := client.ReadJSON(&msg)
		if err != nil {
			h.Logger.ErrorLogging(http.StatusBadRequest, "ReadJSON: "+err.Error())
			return
		}

		msg, err = h.Chat.SendMessage(currentUser, msg)
		if err != nil {
			h.Logger.ErrorLogging(http.StatusBadRequest, "UserUCase: "+err.Error())
			return
		}

		err = client.WriteJSON(msg)
		if err != nil {
			h.Logger.ErrorLogging(http.StatusBadRequest, "WriteJSON")
			return
		}
	}
}

func (h *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	fromId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}
	lastId, err := strconv.Atoi(mux.Vars(r)["lastId"])
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	mses, err := h.Chat.GetChat(r.Context(), uint64(fromId), uint64(lastId))
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, mses)
}

func (h *ChatHandler) GetChats(w http.ResponseWriter, r *http.Request) {
	chats, err := h.Chat.GetChats(r.Context())
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, chats)
}
