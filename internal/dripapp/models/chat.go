package models

import (
	"context"
	"github.com/gorilla/websocket"
	"time"
)

type Message struct {
	MessageID uint64    `json:"messageID" db:"message_id"`
	FromID    uint64    `json:"fromID" db:"from_id"`
	ToID      uint64    `json:"toID" db:"to_id"`
	Text      string    `json:"text"`
	Date      time.Time `json:"date"`
}

type Chat struct {
	FromUserID uint64    `json:"fromUserID"`
	Name       string    `json:"name"`
	Img        string    `json:"img"`
	Messages   []Message `json:"messages"`
}

type ChatUseCase interface {
	CreateClient(c context.Context, conn *websocket.Conn) error
	GetChats(c context.Context) ([]Chat, error)
	GetChat(c context.Context, fromId uint64, lastId uint64) ([]Message, error)
	SaveMessage(currentUser User, message Message) (Message, error)
}

type ChatRepository interface {
	GetChats(ctx context.Context, currentUserId uint64) ([]Chat, error)
	GetChat(ctx context.Context, currentId uint64, fromId uint64, lastId uint64) ([]Message, error)
	SendMessage(currentId uint64, toId uint64, text string) (Message, error)
}
