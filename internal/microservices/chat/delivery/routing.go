package delivery

import (
	_userModels "dripapp/internal/dripapp/models"
	"dripapp/internal/microservices/chat/models"
	"dripapp/internal/pkg/logger"
	_p "dripapp/internal/pkg/permissions"

	"github.com/gorilla/mux"
)

func SetChatRouting(logger logger.Logger, router *mux.Router, chatUCase models.ChatUseCase, ur _userModels.UserRepository) {
	chatHandler := &ChatHandler{
		Chat:   chatUCase,
		Logger: logger,
	}

	userMid := _p.UserMiddlware{
		UserRepo: ur,
	}

	router.HandleFunc("/api/v1/notifications",
		_p.SetCSRF(_p.CheckAuthenticated(userMid.GetCurrentUser(chatHandler.Notifications))))

	router.HandleFunc("/api/v1/chat/{id:[0-9]+}&{lastId:[0-9]+}",
		_p.SetCSRF(_p.CheckAuthenticated(userMid.GetCurrentUser(chatHandler.GetChat)))).
		Methods("GET", "OPTIONS")

	router.HandleFunc("/api/v1/chats",
		_p.SetCSRF(_p.CheckAuthenticated(userMid.GetCurrentUser(chatHandler.GetChats)))).
		Methods("GET", "OPTIONS")
}
