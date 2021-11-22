package usecase

import (
	"context"
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	"github.com/gorilla/websocket"
	"time"
)

type ChatUseCase struct {
	ChatRepo       models.ChatRepository
	hub            *Hub
	contextTimeout time.Duration
}

func NewChatUseCase(
	chatRepo models.ChatRepository,
	hub *Hub,
	timeout time.Duration) models.ChatUseCase {
	return &ChatUseCase{
		ChatRepo:       chatRepo,
		hub:            hub,
		contextTimeout: timeout,
	}
}

func (h *ChatUseCase) CreateClient(c context.Context, conn *websocket.Conn) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return models.ErrContextNilError
	}

	client := &Client{
		user: currentUser,
		uc:   h,
		conn: conn,
		send: make(chan models.Message),
	}
	h.hub.register <- client

	go client.writePump()
	go client.readPump()

	return nil
}

func (h *ChatUseCase) GetChats(c context.Context) ([]models.Chat, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return nil, models.ErrContextNilError
	}

	chats, err := h.ChatRepo.GetChats(ctx, currentUser.ID)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (h *ChatUseCase) GetChat(c context.Context, fromId uint64, lastId uint64) ([]models.Message, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return nil, models.ErrContextNilError
	}

	mses, err := h.ChatRepo.GetChat(ctx, currentUser.ID, fromId, lastId)
	if err != nil {
		return nil, err
	}

	return mses, nil
}

func (h *ChatUseCase) SaveMessage(currentUser models.User, message models.Message) (models.Message, error) {
	msg, err := h.ChatRepo.SendMessage(currentUser.ID, message.ToID, message.Text)
	if err != nil {
		return models.Message{}, err
	}

	return msg, nil
}
