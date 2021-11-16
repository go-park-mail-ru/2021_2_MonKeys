package delivery

import (
	"bytes"
	"context"
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	_s "dripapp/internal/dripapp/session/mocks"
	"dripapp/internal/dripapp/user/mocks"
	"dripapp/internal/pkg/logger"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

type TestCase struct {
	BodyReq         io.Reader
	StatusCode      int
	BodyResp        string
	mockUserUseCase []interface{}
	mockSessUseCase []interface{}
	SessionCookie   http.Cookie
}

func TestLogin(t *testing.T) {
	t.Parallel()

	mockUserUseCase := &mocks.UserUsecase{}
	mockSessionUseCase := &_s.SessionUsecase{}

	sessionHandler := &SessionHandler{
		Logger:       logger.DripLogger,
		UserUCase:    mockUserUseCase,
		SessionUcase: mockSessionUseCase,
	}

	cases := []TestCase{
		{
			BodyReq:    bytes.NewReader([]byte(`{"email":"testLogin1@mail.ru","password":"123456qQ"}`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":{"id":1,"email":"testLogin1@mail.ru"}}`,
			mockUserUseCase: []interface{}{
				models.User{
					ID:       1,
					Email:    "testLogin1@mail.ru",
					Password: "123456qQ",
				},
				nil,
			},
			mockSessUseCase: []interface{}{
				nil,
			},
		},
		{
			BodyReq:    bytes.NewReader([]byte(`wrong input data`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":400,"body":null}`,
			mockUserUseCase: []interface{}{
				models.User{},
				errors.New(""),
			},
			mockSessUseCase: []interface{}{
				nil,
			},
		},
		{
			BodyReq:    bytes.NewReader([]byte(`{"email":"wrongEmail","password":"wrongPassword"}`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
			mockUserUseCase: []interface{}{
				models.User{},
				errors.New(""),
			},
			mockSessUseCase: []interface{}{
				nil,
			},
		},
		{
			BodyReq:    bytes.NewReader([]byte(`{"email":"testLogin1@mail.ru","password":"wrongPassword"}`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
			mockUserUseCase: []interface{}{
				models.User{},
				errors.New(""),
			},
			mockSessUseCase: []interface{}{
				nil,
			},
		},
		{
			BodyReq:    bytes.NewReader([]byte(`{"email":"testLogin1@mail.ru","password":"123456qQ"}`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":500,"body":null}`,
			mockUserUseCase: []interface{}{
				models.User{
					ID:       1,
					Email:    "testLogin1@mail.ru",
					Password: "123456qQ",
				},
				nil,
			},
			mockSessUseCase: []interface{}{
				errors.New("session already exists"),
			},
		},
	}

	for caseNum, item := range cases {
		r := httptest.NewRequest("POST", "/api/v1/login", item.BodyReq)
		r = r.WithContext(context.WithValue(r.Context(), configs.ContextUserID, models.Session{
			UserID: 0,
			Cookie: "",
		}))
		w := httptest.NewRecorder()

		mockUserUseCase.ExpectedCalls = nil
		mockUserUseCase.On("Login",
			r.Context(),
			mock.AnythingOfType("models.LoginUser")).Return(item.mockUserUseCase...)
		mockSessionUseCase.ExpectedCalls = nil
		mockSessionUseCase.On("AddSession",
			r.Context(),
			mock.AnythingOfType("models.Session")).Return(item.mockSessUseCase...)

		sessionHandler.LoginHandler(w, r)

		if w.Code != item.StatusCode {
			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
				caseNum, w.Code, item.StatusCode)
		}

		if w.Body.String() != item.BodyResp {
			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
				caseNum, w.Body.String(), item.BodyResp)
		}
	}
}

func TestLogout(t *testing.T) {
	t.Parallel()

	mockUserUseCase := &mocks.UserUsecase{}
	mockSessionUseCase := &_s.SessionUsecase{}

	sessionHandler := &SessionHandler{
		Logger:       logger.DripLogger,
		UserUCase:    mockUserUseCase,
		SessionUcase: mockSessionUseCase,
	}

	cases := []TestCase{
		{
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":null}`,
			mockSessUseCase: []interface{}{
				nil,
			},
			SessionCookie: http.Cookie{
				Name: "sessionId",
			},
		},
		{
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
			mockSessUseCase: []interface{}{
				errors.New("session does not exist"),
			},
			SessionCookie: http.Cookie{
				Name: "sessionId",
			},
		},
	}

	for caseNum, item := range cases {
		r := httptest.NewRequest("GET", "/api/v1/logout", item.BodyReq)
		r.AddCookie(&item.SessionCookie)
		r = r.WithContext(context.WithValue(r.Context(), configs.ContextUserID, models.Session{
			UserID: 0,
			Cookie: "",
		}))
		w := httptest.NewRecorder()

		mockSessionUseCase.ExpectedCalls = nil
		mockSessionUseCase.On("DeleteSession", r.Context()).Return(item.mockSessUseCase...)

		sessionHandler.LogoutHandler(w, r)

		if w.Code != item.StatusCode {
			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
				caseNum, w.Code, item.StatusCode)
		}

		if w.Body.String() != item.BodyResp {
			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
				caseNum, w.Body.String(), item.BodyResp)
		}
	}
}

func TestSetRouting(t *testing.T) {
	mockUserUseCase := &mocks.UserUsecase{}
	mockSessionUseCase := &_s.SessionUsecase{}

	SetSessionRouting(logger.DripLogger, mux.NewRouter(), mockUserUseCase, mockSessionUseCase)
}
