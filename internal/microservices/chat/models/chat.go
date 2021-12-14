package models

import (
	"context"
	"time"
)

//easyjson:json
type Message struct {
	MessageID uint64    `json:"messageID" db:"message_id"`
	FromID    uint64    `json:"fromID" db:"from_id"`
	ToID      uint64    `json:"toID" db:"to_id"`
	Text      string    `json:"text"`
	Date      time.Time `json:"date"`
}

type Messages struct {
	Messages []Message
}

type Chat struct {
	FromUserID uint64    `json:"fromUserID" db:"id"`
	Name       string    `json:"name"`
	Img        string    `json:"img"`
	Messages   []Message `json:"messages"`
}

type Chats struct {
	Chats []Chat
}

type ChatUseCase interface {
	ClientHandler(c context.Context, io IOMessage) error
	GetChats(c context.Context) ([]Chat, error)
	GetChat(c context.Context, fromId uint64, lastId uint64) ([]Message, error)
	DeleteChat(c context.Context, fromId uint64) error
}

type ChatRepository interface {
	GetChats(ctx context.Context, userId uint64) ([]Chat, error)
	GetChat(ctx context.Context, userId uint64, fromId uint64, lastId uint64) ([]Message, error)
	SaveMessage(userId uint64, toId uint64, text string) (Message, error)
	DeleteChat(ctx context.Context, userId uint64, fromId uint64) error
}
