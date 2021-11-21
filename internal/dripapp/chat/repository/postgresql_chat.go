package repository

import (
	"context"
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type PostgreChatRepo struct {
	Conn sqlx.DB
}

const success = "Connection success (postgre) on: "

func NewPostgresChatRepository(config configs.PostgresConfig) (models.ChatRepository, error) {
	ConnStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable",
		config.User,
		config.DBName,
		config.Password,
		config.Host)

	Conn, err := sqlx.Open("postgres", ConnStr)
	if err != nil {
		return nil, err
	}

	log.Printf("%s%s", success, ConnStr)
	return &PostgreChatRepo{*Conn}, nil
}

func (p PostgreChatRepo) GetChats(ctx context.Context, currentUserId uint64) ([]models.Chat, error) {
	var chats []models.Chat
	err := p.Conn.Select(&chats, GetChats, currentUserId)
	if err != nil {
		return nil, err
	}

	for idx := range chats {
		var ms models.Message
		err := p.Conn.GetContext(ctx, &ms, GetLastMessage, currentUserId, chats[idx].FromUserID)
		if err != nil {
			return nil, err
		}
		chats[idx].Messages = append(chats[idx].Messages, ms)
	}

	return chats, nil
}

func (p PostgreChatRepo) GetChat(ctx context.Context, currentId uint64, fromId uint64, lastId uint64) ([]models.Message, error) {
	var mses []models.Message
	err := p.Conn.Select(&mses, GetMessages, currentId, fromId, lastId)
	if err != nil {
		return nil, err
	}

	return mses, nil
}

func (p PostgreChatRepo) SendMessage(currentId uint64, toId uint64, text string) (models.Message, error) {
	var msg models.Message
	err := p.Conn.Get(&msg, SendNessage, currentId, toId, text)
	if err != nil {
		return models.Message{}, err
	}

	return msg, nil
}
