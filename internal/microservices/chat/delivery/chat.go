package delivery

import (
	_userModels "dripapp/internal/dripapp/models"
	"dripapp/internal/microservices/chat/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/responses"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ChatHandler struct {
	Chat   models.ChatUseCase
	Logger logger.Logger
}

func (h *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	fromId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		responses.SendError(w, _userModels.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}
	lastId, err := strconv.Atoi(mux.Vars(r)["lastId"])
	if err != nil {
		responses.SendError(w, _userModels.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	mses, err := h.Chat.GetChat(r.Context(), uint64(fromId), uint64(lastId))
	if err != nil {
		responses.SendError(w, _userModels.HTTPError{
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
		responses.SendError(w, _userModels.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, chats)
}
