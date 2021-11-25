package usecase

import (
	"context"
	"dripapp/configs"
	_userModels "dripapp/internal/dripapp/models"
	_chatModels "dripapp/internal/microservices/chat/models"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	chatMocks "dripapp/internal/microservices/chat/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChatUsecase_GetChats(t *testing.T) {
	type TestCase struct {
		user  _userModels.User
		chats []_chatModels.Chat
		err   error
	}
	testCases := []TestCase{
		// Test OK
		{
			user: _userModels.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			chats: []_chatModels.Chat{
				{
					FromUserID: 1,
					Name:       "Vova",
					Img:        "/media/test.webp",
					Messages:   []_chatModels.Message{},
				},
			},
			err: nil,
		},
		// Test ErrContextNilError
		{
			user: _userModels.User{
				ID:          2,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			chats: []_chatModels.Chat{},
			err:   _userModels.ErrContextNilError,
		},
		// Test GetChatsErr
		{
			user: _userModels.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			chats: []_chatModels.Chat{},
			err:   errors.New(""),
		},
	}

	type MockResultCase struct {
		chats []_chatModels.Chat
		err   error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			chats: []_chatModels.Chat{
				{
					FromUserID: 1,
					Name:       "Vova",
					Img:        "/media/test.webp",
					Messages:   []_chatModels.Message{},
				},
			},
			err: nil,
		},
		// Test ErrContextNilError
		{
			chats: []_chatModels.Chat{},
			err:   nil,
		},
		// Test GetChatsErr
		{
			chats: []_chatModels.Chat{},
			err:   errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.user.ID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, testCase.user))
		}

		mockSessionRepository := new(chatMocks.ChatRepository)
		mockSessionRepository.On("GetChats",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].chats, MockResultCases[i].err)

		testHub := _chatModels.NewHub()
		testChatUsecase := NewChatUseCase(mockSessionRepository, testHub, time.Second*2)
		resultChats, err := testChatUsecase.GetChats(r.Context())

		assert.Equal(t, testCase.err, err, testCase.err, message)
		reflect.DeepEqual(resultChats, testCase.chats)
	}
}

func TestChatUsecase_GetChat(t *testing.T) {
	type TestCase struct {
		user     _userModels.User
		fromId   uint64
		lastId   uint64
		messages []_chatModels.Message
		err      error
	}
	testCases := []TestCase{
		// Test OK
		{
			user: _userModels.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			messages: []_chatModels.Message{
				{
					MessageID: 1,
					FromID:    0,
					ToID:      1,
					Text:      "Hello",
				},
				{
					MessageID: 2,
					FromID:    1,
					ToID:      0,
					Text:      "Hello!",
				},
			},
			err: nil,
		},
		// Test ErrContextNilError
		{
			user: _userModels.User{
				ID:          2,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			messages: []_chatModels.Message{},
			err:      _userModels.ErrContextNilError,
		},
		// Test GetChatsErr
		{
			user: _userModels.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			messages: []_chatModels.Message{},
			err:      errors.New(""),
		},
	}

	type MockResultCase struct {
		chat []_chatModels.Message
		err  error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			chat: []_chatModels.Message{
				{
					MessageID: 1,
					FromID:    0,
					ToID:      1,
					Text:      "Hello",
				},
				{
					MessageID: 2,
					FromID:    1,
					ToID:      0,
					Text:      "Hello!",
				},
			},
			err: nil,
		},
		// Test ErrContextNilError
		{
			chat: []_chatModels.Message{},
			err:  nil,
		},
		// Test GetChatsErr
		{
			chat: []_chatModels.Message{},
			err:  errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.user.ID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, testCase.user))
		}

		mockSessionRepository := new(chatMocks.ChatRepository)
		mockSessionRepository.On("GetChat",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].chat, MockResultCases[i].err)

		testHub := _chatModels.NewHub()
		testChatUsecase := NewChatUseCase(mockSessionRepository, testHub, time.Second*2)
		resultChats, err := testChatUsecase.GetChat(r.Context(), testCase.fromId, testCase.lastId)

		assert.Equal(t, testCase.err, err, testCase.err, message)
		reflect.DeepEqual(resultChats, testCase.messages)
	}
}
