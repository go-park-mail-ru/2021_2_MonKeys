package usecase_test

import (
	"context"
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	sessionMocks "dripapp/internal/dripapp/session/mocks"
	"dripapp/internal/dripapp/session/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func CreateMultipartRequest(method, target string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, target, body)
	if err != nil {
		return nil, err
	}
	r = r.WithContext(context.WithValue(r.Context(), configs.ForContext, models.Session{
		UserID: 0,
		Cookie: "",
	}))

	return r, nil
}

func TestUserUsecase_AddSession(t *testing.T) {
	type TestCase struct {
		session models.Session
		err     error
	}
	testCases := []TestCase{
		// Test OK
		{
			session: models.Session{
				Cookie: "",
				UserID: 0,
			},
			err: nil,
		},
		// Test Err
		{
			session: models.Session{
				Cookie: "",
				UserID: 0,
			},
			err: errors.New(""),
		},
	}

	type MockResultCase struct {
		err error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			err: nil,
		},
		// Test Err
		{
			err: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		r = r.WithContext(context.WithValue(r.Context(), configs.ForContext, testCase.session))

		mockSessionRepository := new(sessionMocks.SessionRepository)
		mockSessionRepository.On("NewSessionCookie", mock.AnythingOfType("string"), mock.AnythingOfType("uint64")).Return(MockResultCases[i].err)

		testSessionUsecase := usecase.NewSessionUsecase(mockSessionRepository, time.Second*2)
		err = testSessionUsecase.AddSession(r.Context(), testCase.session)

		assert.Equal(t, testCase.err, err, testCase.err, message)
	}
}

func TestUserUsecase_DeleteSession(t *testing.T) {
	type TestCase struct {
		session models.Session
		err     error
	}
	testCases := []TestCase{
		// Test OK
		{
			session: models.Session{
				Cookie: "",
				UserID: 0,
			},
			err: nil,
		},
		// Test Err
		{
			session: models.Session{
				Cookie: "",
				UserID: 1,
			},
			err: errors.New("context nil error"),
		},
		// Test ErrDeleteSessionCookie
		{
			session: models.Session{
				Cookie: "",
				UserID: 0,
			},
			err: errors.New(""),
		},
	}

	type MockResultCase struct {
		err error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			err: nil,
		},
		// Test ErrContextNil
		{
			err: nil,
		},
		// Test ErrDeleteSessionCookie
		{
			err: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.session.UserID != 1 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ForContext, testCase.session))
		}

		mockSessionRepository := new(sessionMocks.SessionRepository)
		mockSessionRepository.On("DeleteSessionCookie", mock.AnythingOfType("string")).Return(MockResultCases[i].err)

		testSessionUsecase := usecase.NewSessionUsecase(mockSessionRepository, time.Second*2)
		err = testSessionUsecase.DeleteSession(r.Context())

		assert.Equal(t, testCase.err, err, testCase.err, message)
	}
}
