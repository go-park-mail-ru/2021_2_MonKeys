package delivery

import (
	_authClient "dripapp/internal/microservices/auth/delivery/grpc/client"
	"dripapp/internal/microservices/chat/models"
	"dripapp/internal/pkg/logger"
	_p "dripapp/internal/pkg/permissions"

	"github.com/gorilla/mux"
)

func SetChatRouting(logger logger.Logger, router *mux.Router, chatUCase models.ChatUseCase, sc _authClient.SessionClient) {
	chatHandler := &ChatHandler{
		Chat:   chatUCase,
		Logger: logger,
	}

	perm := _p.Permission{
		AuthClient: sc,
	}

	router.HandleFunc("/api/v1/apiws",
		_p.SetCSRF(perm.CheckAuth(perm.GetCurrentUser(chatHandler.UpgradeWS))))

	router.HandleFunc("/api/v1/chat/{id:[0-9]+}&{lastId:[0-9]+}",
		_p.SetCSRF(perm.CheckAuth(perm.GetCurrentUser(chatHandler.GetChat)))).
		Methods("GET", "OPTIONS")

	router.HandleFunc("/api/v1/chats",
		_p.SetCSRF(perm.CheckAuth(perm.GetCurrentUser(chatHandler.GetChats)))).
		Methods("GET", "OPTIONS")
}
