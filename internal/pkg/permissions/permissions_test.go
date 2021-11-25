package permissions

import (
	"context"
	"dripapp/configs"
	_userModels "dripapp/internal/dripapp/models"
	_authClient "dripapp/internal/microservices/auth/delivery/grpc/client"
	_authMock "dripapp/internal/microservices/auth/mocks"
	_sessionModels "dripapp/internal/microservices/auth/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestSetCsrf(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrf := w.Header().Get("csrf")

		assert.NotEqual(t, "", csrf)
	})
	handlerToTest := SetCSRF(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)

	recorder := httptest.NewRecorder()
	handlerToTest.ServeHTTP(recorder, req)

}

func TestCheckCsrf(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrf := w.Header().Get("csrf")
		fmt.Println(csrf)
		assert.Equal(t, "", csrf)
	})
	handlerToTest := CheckCSRF(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)

	recorder := httptest.NewRecorder()
	handlerToTest.ServeHTTP(recorder, req)
}

func TestCheckAuth(t *testing.T) {
	t.Parallel()
	grpcConn, _ := grpc.Dial(configs.AuthServer.GrpcUrl, grpc.WithInsecure())
	grpcAuthClient := _authClient.NewAuthClient(grpcConn)

	perm := Permission{
		AuthClient: *grpcAuthClient,
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, true)
	})
	handlerToTest := perm.CheckAuth(nextHandler)

	t.Run("good check auth", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:   "sessionId",
			Value:  "qwerty",
			MaxAge: 300,
		}

		req := httptest.NewRequest("GET", "http://testing", nil)
		req.AddCookie(cookie)

		m := new(_authMock.AuthGrpcHandlerClient)
		m.On("GetFromSession", req.Context(), "qwerty").Return(_sessionModels.Session{}, nil)

		recorder := httptest.NewRecorder()
		handlerToTest.ServeHTTP(recorder, req)
	})
}

func TestGetCurrentUser(t *testing.T) {
	grpcConn, _ := grpc.Dial(configs.AuthServer.GrpcUrl, grpc.WithInsecure())
	grpcAuthClient := _authClient.NewAuthClient(grpcConn)

	perm := Permission{
		AuthClient: *grpcAuthClient,
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, true)
	})
	handlerToTest := perm.CheckAuth(nextHandler)

	t.Run("good get current", func(t *testing.T) {

		sess := _sessionModels.Session{
			Cookie: "qwerty",
			UserID: 1,
		}

		req := httptest.NewRequest("GET", "http://testing", nil)

		req = req.WithContext(context.WithValue(req.Context(), configs.ContextUserID, 1))

		m := new(_authMock.AuthGrpcHandlerClient)
		m.On("GetFromSession", req.Context(), sess).Return(_userModels.User{}, nil)

		recorder := httptest.NewRecorder()
		handlerToTest.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("TestCase [%s]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
				"good get current", recorder.Code, http.StatusOK)
		}
		if recorder.Body.String() != `{"status":403,"body":null}` {
			t.Errorf("TestCase [%s]:\nwrongCase Response: \ngot %s\nexpected %s",
				"good get current", recorder.Body.String(), `{"status":403,"body":null}`)
		}
	})
}
