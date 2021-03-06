package delivery

import (
	"context"
	"dripapp/configs"
	_models "dripapp/internal/dripapp/models"
	_authClient "dripapp/internal/microservices/auth/delivery/grpc/client"
	"dripapp/internal/microservices/chat/mocks"
	"dripapp/internal/microservices/chat/models"
	"dripapp/internal/pkg/logger"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type TestCase struct {
	URL        string
	BodyReq    io.Reader
	MockChat   []interface{}
	StatusCode int
	BodyResp   string
}

var (
	message1 = models.Message{
		MessageID: 1,
		FromID:    1,
		ToID:      2,
		Text:      "text from 1",
		Date:      time.Now(),
	}
	message2 = models.Message{
		MessageID: 2,
		FromID:    2,
		ToID:      1,
		Text:      "text from 2",
		Date:      time.Now(),
	}
	messages = []models.Message{
		message1,
		message2,
	}
	messagesEasy = models.Messages{
		Messages: messages,
	}
	messagesStr = objToJsonStr(messagesEasy)

	chat1 = models.Chat{
		FromUserID: 1,
		Name:       "chat name",
		Img:        "chat.img",
		Messages:   messages,
	}
	chat2 = models.Chat{
		FromUserID: 2,
		Name:       "chat name",
		Img:        "chat.img",
		Messages:   messages,
	}
	chats = []models.Chat{
		chat1,
		chat2,
	}
	chatsEasy = models.Chats{
		Chats: chats,
	}
	chatsStr = objToJsonStr(chatsEasy)
)

func objToJsonStr(v interface{}) string {
	j, _ := json.Marshal(v)
	return string(j)
}

func CheckResponse(t *testing.T, w *httptest.ResponseRecorder, caseNum int, testCase TestCase) {
	if w.Code != testCase.StatusCode {
		t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
			caseNum, w.Code, testCase.StatusCode)
	}

	if w.Body.String() != testCase.BodyResp {
		t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
			caseNum, w.Body.String(), testCase.BodyResp)
	}
}

func CreateRequest(method, target string, body io.Reader) (r *http.Request) {
	r = httptest.NewRequest(method, target, body)
	r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, _models.User{}))

	return
}

func TestSetRouting(t *testing.T) {
	mockChat := &mocks.ChatUseCase{}

	grpcConn, _ := grpc.Dial(configs.AuthServer.GrpcUrl, grpc.WithInsecure())
	grpcAuthClient := _authClient.NewAuthClient(grpcConn)

	SetChatRouting(logger.DripLogger, mux.NewRouter(), mockChat, *grpcAuthClient)
}

func TestGetChat(t *testing.T) {
	t.Parallel()

	mockChat := &mocks.ChatUseCase{}

	chatHandler := &ChatHandler{
		Chat:   mockChat,
		Logger: logger.DripLogger,
	}

	cases := []TestCase{
		{
			MockChat: []interface{}{
				messages,
				nil,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":` + messagesStr + `}`,
		},
		{
			MockChat: []interface{}{
				messages,
				errors.New(""),
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
	}

	for caseNum, item := range cases {
		r := httptest.NewRequest("POST", "/api/v1/chat/1&1", item.BodyReq)
		vars := map[string]string{
			"id":     "1",
			"lastId": "1",
		}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()

		mockChat.ExpectedCalls = nil
		mockChat.On("GetChat",
			r.Context(),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("uint64")).Return(item.MockChat...)

		chatHandler.GetChat(w, r)

		CheckResponse(t, w, caseNum, item)
	}
}

func TestGetChats(t *testing.T) {
	t.Parallel()

	mockChat := &mocks.ChatUseCase{}

	chatHandler := &ChatHandler{
		Chat:   mockChat,
		Logger: logger.DripLogger,
	}

	call := mockChat.On("GetChats", context.Background())

	cases := []TestCase{
		{
			MockChat: []interface{}{
				chats,
				nil,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":` + chatsStr + `}`,
		},
		{
			MockChat: []interface{}{
				chats,
				errors.New(""),
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
	}

	for caseNum, item := range cases {
		call.Return(item.MockChat...)

		r := httptest.NewRequest("GET", "/api/v1/chats", nil)
		w := httptest.NewRecorder()

		chatHandler.GetChats(w, r)

		CheckResponse(t, w, caseNum, item)
	}
}

func TestDeleteChat(t *testing.T) {
	t.Parallel()

	mockChat := &mocks.ChatUseCase{}

	chatHandler := &ChatHandler{
		Chat:   mockChat,
		Logger: logger.DripLogger,
	}

	cases := []TestCase{
		{
			URL: "1",
			MockChat: []interface{}{
				nil,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":null}`,
		},
		{
			URL: "wrongURL",
			MockChat: []interface{}{
				nil,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
		{
			URL: "1",
			MockChat: []interface{}{
				errors.New(""),
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
	}

	for caseNum, item := range cases {
		r := CreateRequest("DELETE", "/api/v1/chat/"+item.URL, nil)
		vars := map[string]string{
			"id": item.URL,
		}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()

		mockChat.ExpectedCalls = nil
		mockChat.On("DeleteChat",
			r.Context(),
			mock.AnythingOfType("uint64")).Return(item.MockChat...)

		chatHandler.DeleteChat(w, r)

		CheckResponse(t, w, caseNum, item)
	}
}
