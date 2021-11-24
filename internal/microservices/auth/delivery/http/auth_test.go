package http

import (
	"bytes"
	"context"
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	"dripapp/internal/dripapp/user/mocks"
	_authClient "dripapp/internal/microservices/auth/delivery/grpc/client"
	_s "dripapp/internal/microservices/auth/mocks"
	_sessionModels "dripapp/internal/microservices/auth/models"
	"dripapp/internal/pkg/logger"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type TestCase struct {
	BodyReq         io.Reader
	StatusCode      int
	BodyResp        string
	mockUserUseCase []interface{}
	mockSessUseCase []interface{}
	SessionCookie   http.Cookie
}

var (
	idStr    = "1"
	user     = models.User{
		ID:       uint64(1),
		Email:    "test@mail.ru",
		Password: "qweQWE12",
	}
)

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
	r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, models.User{}))

	return
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
		r = r.WithContext(context.WithValue(r.Context(), configs.ContextUserID, _sessionModels.Session{
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
		r = r.WithContext(context.WithValue(r.Context(), configs.ContextUserID, _sessionModels.Session{
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

func TestSignup(t *testing.T) {
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
			BodyReq: bytes.NewReader([]byte(`{"email":"` + user.Email + `","password":"` + user.Password + `"}`)),
			mockUserUseCase: []interface{}{
				user,
				nil,
			},
			mockSessUseCase: []interface{}{
				nil,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":{"id":` + idStr + `,"email":"` + user.Email + `"}}`,
		},
		{
			BodyReq: bytes.NewReader([]byte(`wrong input data`)),
			mockUserUseCase: []interface{}{
				user,
				nil,
			},
			mockSessUseCase: []interface{}{
				nil,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":400,"body":null}`,
		},
		{
			BodyReq: bytes.NewReader([]byte(`{"email":"wrongEmail","password":"wrongPassword"}`)),
			mockUserUseCase: []interface{}{
				models.User{},
				models.ErrEmailAlreadyExists,
			},
			mockSessUseCase: []interface{}{
				nil,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":1001,"body":null}`,
		},
		{
			BodyReq: bytes.NewReader([]byte(`{"email":"` + user.Email + `","password":"` + user.Password + `"}`)),
			mockUserUseCase: []interface{}{
				user,
				nil,
			},
			mockSessUseCase: []interface{}{
				models.ErrSessionAlreadyExists,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":500,"body":null}`,
		},
	}

	for caseNum, item := range cases {
		r := CreateRequest("POST", "/api/v1/signup", item.BodyReq)
		w := httptest.NewRecorder()

		mockUserUseCase.ExpectedCalls = nil
		mockUserUseCase.On("Signup",
			r.Context(),
			mock.AnythingOfType("models.LoginUser")).Return(item.mockUserUseCase...)
		mockSessionUseCase.ExpectedCalls = nil
		mockSessionUseCase.On("AddSession",
			r.Context(),
			mock.AnythingOfType("models.Session")).Return(item.mockSessUseCase...)

		sessionHandler.SignupHandler(w, r)

		CheckResponse(t, w, caseNum, item)
	}
}

func TestSetRouting(t *testing.T) {
	mockUserUseCase := &mocks.UserUsecase{}
	mockSessionUseCase := &_s.SessionUsecase{}
	grpcConn, _ := grpc.Dial(configs.AuthServer.GrpcUrl, grpc.WithInsecure())
	grpcAuthClient := _authClient.NewStaffClient(grpcConn)

	SetSessionRouting(logger.DripLogger, mux.NewRouter(), mockUserUseCase, mockSessionUseCase, *grpcAuthClient)
}
