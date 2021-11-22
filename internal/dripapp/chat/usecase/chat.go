package usecase

import (
	"context"
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	"time"
)

type ChatUseCase struct {
	ChatRepo       models.ChatRepository
	hub            *models.Hub
	contextTimeout time.Duration
}

func NewChatUseCase(
	chatRepo models.ChatRepository,
	hub *models.Hub,
	timeout time.Duration) models.ChatUseCase {
	return &ChatUseCase{
		ChatRepo:       chatRepo,
		hub:            hub,
		contextTimeout: timeout,
	}
}

func (h *ChatUseCase) ClientHandler(c context.Context, io models.IOMessage) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return models.ErrContextNilError
	}

	client := models.NewChatClient(currentUser, h.hub, h.ChatRepo, io)

	go client.WritePump()
	go client.ReadPump()

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
