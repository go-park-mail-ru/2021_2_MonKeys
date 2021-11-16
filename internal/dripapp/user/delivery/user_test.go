package delivery

import (
	"bytes"
	"context"
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	_s "dripapp/internal/dripapp/session/mocks"
	"dripapp/internal/dripapp/user/mocks"
	"dripapp/internal/dripapp/user/repository"
	"dripapp/internal/pkg/logger"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

// const (
// 	correctCase = iota + 1
// 	wrongCase
// )

type TestCase struct {
	BodyReq         io.Reader
	mockUserUseCase []interface{}
	mockSessUseCase []interface{}
	StatusCode      int
	BodyResp        string
}

var (
	email    = "test@mail.ru"
	password = "qweQWE12"
	id       = uint64(1)
	idStr    = "1"
	user     = models.User{
		ID:       id,
		Email:    email,
		Password: password,
	}
	user2 = models.User{
		ID:       2,
		Email:    "test2@mail.ru",
		Password: "qweQWE12",
	}

	tags = models.Tags{
		AllTags: map[uint64]models.Tag{
			1: {Tag_Name: "chill"},
			2: {Tag_Name: "sport"},
			3: {Tag_Name: "music"},
		},
		Count: 3,
	}
	tagsMapStr   = `{"1":{"tagText":"chill"},"2":{"tagText":"sport"},"3":{"tagText":"music"}}`
	tagsCountStr = "3"

	matches = models.Matches{
		AllUsers: map[uint64]models.User{
			1: user,
			2: user2,
		},
		Count: "2",
	}
	usersMapStr     = `{"1":{"id":1,"email":"test@mail.ru"},"2":{"id":2,"email":"test2@mail.ru"}}`
	matchesCountStr = "2"

	reactionStr = "0"
	match       = models.Match{Match: true}
	matchStr    = "true"

	notMatch    = models.Match{Match: false}
	notMatchStr = "false"

	photo = models.Photo{Path: "path"}
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
	r = r.WithContext(context.WithValue(r.Context(), configs.ContextUserID, models.Session{
		UserID: 0,
		Cookie: "",
	}))

	return
}

func TestCurrentUser(t *testing.T) {
	t.Parallel()

	mockUserUseCase := &mocks.UserUsecase{}
	mockSessionUseCase := &_s.SessionUsecase{}

	call := mockUserUseCase.On("CurrentUser", context.Background())

	userHandler := &UserHandler{
		Logger:       logger.DripLogger,
		UserUCase:    mockUserUseCase,
		SessionUcase: mockSessionUseCase,
	}

	cases := []TestCase{
		{
			mockUserUseCase: []interface{}{
				user,
				models.StatusOk200,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":{"id":` + idStr + `,"email":"` + email + `"}}`,
		},
		{
			mockUserUseCase: []interface{}{
				models.User{},
				models.HTTPError{
					Code:    http.StatusNotFound,
					Message: models.ErrContextNilError,
				},
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
	}

	for caseNum, item := range cases {
		call.Return(item.mockUserUseCase...)

		r := httptest.NewRequest("GET", "/api/v1/currentuser", nil)
		w := httptest.NewRecorder()

		userHandler.CurrentUser(w, r)

		CheckResponse(t, w, caseNum, item)
	}
}

func TestSignup(t *testing.T) {
	t.Parallel()

	mockUserUseCase := &mocks.UserUsecase{}
	mockSessionUseCase := &_s.SessionUsecase{}

	userHandler := &UserHandler{
		Logger:       logger.DripLogger,
		UserUCase:    mockUserUseCase,
		SessionUcase: mockSessionUseCase,
	}

	cases := []TestCase{
		{
			BodyReq: bytes.NewReader([]byte(`{"email":"` + email + `","password":"` + password + `"}`)),
			mockUserUseCase: []interface{}{
				user,
				models.StatusOk200,
			},
			mockSessUseCase: []interface{}{
				nil,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":{"id":` + idStr + `,"email":"` + email + `"}}`,
		},
		{
			BodyReq: bytes.NewReader([]byte(`wrong input data`)),
			mockUserUseCase: []interface{}{
				user,
				models.StatusOk200,
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
				models.HTTPError{
					Code: models.StatusEmailAlreadyExists,
				},
			},
			mockSessUseCase: []interface{}{
				nil,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":1001,"body":null}`,
		},
		{
			BodyReq: bytes.NewReader([]byte(`{"email":"` + email + `","password":"` + password + `"}`)),
			mockUserUseCase: []interface{}{
				user,
				models.StatusOk200,
			},
			mockSessUseCase: []interface{}{
				errors.New("session already exists"),
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

		userHandler.SignupHandler(w, r)

		CheckResponse(t, w, caseNum, item)
	}
}

func TestNextUser(t *testing.T) {
	t.Parallel()

	mockUserUseCase := &mocks.UserUsecase{}
	mockSessionUseCase := &_s.SessionUsecase{}

	userHandler := &UserHandler{
		Logger:       logger.DripLogger,
		UserUCase:    mockUserUseCase,
		SessionUcase: mockSessionUseCase,
	}

	cases := []TestCase{
		{
			mockUserUseCase: []interface{}{
				[]models.User{user},
				models.StatusOk200,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":[{"id":` + idStr + `,"email":"` + email + `"}]}`,
		},
		{
			mockUserUseCase: []interface{}{
				[]models.User{},
				models.HTTPError{Code: http.StatusNotFound},
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
	}

	for caseNum, item := range cases {
		r := CreateRequest("GET", "/api/v1/nextswipeuser", nil)
		w := httptest.NewRecorder()

		mockUserUseCase.ExpectedCalls = nil
		mockUserUseCase.On("NextUser", r.Context()).Return(item.mockUserUseCase...)

		userHandler.NextUserHandler(w, r)

		CheckResponse(t, w, caseNum, item)
	}
}

func TestEditProfile(t *testing.T) {
	t.Parallel()

	mockUserUseCase := &mocks.UserUsecase{}
	mockSessionUseCase := &_s.SessionUsecase{}

	userHandler := &UserHandler{
		Logger:       logger.DripLogger,
		UserUCase:    mockUserUseCase,
		SessionUcase: mockSessionUseCase,
	}

	cases := []TestCase{
		{
			BodyReq: bytes.NewReader([]byte(`{"email":"` + email + `","password":"` + password + `"}`)),
			mockUserUseCase: []interface{}{
				user,
				models.StatusOk200,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":{"id":` + idStr + `,"email":"` + email + `"}}`,
		},
		{
			BodyReq:    bytes.NewReader([]byte(`wrong input data`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":400,"body":null}`,
		},
		{
			BodyReq: bytes.NewReader([]byte(`{"name":"testEdit","date":"wrong-format-data","description":"Description Description Description Description","imgSrc":"/img/testEdit/","tags":["Tags","Tags","Tags","Tags","Tags"]}`)),
			mockUserUseCase: []interface{}{
				models.User{},
				models.HTTPError{Code: http.StatusBadRequest},
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":400,"body":null}`,
		},
	}

	for caseNum, item := range cases {
		r := CreateRequest("POST", "/api/v1/editprofile", item.BodyReq)
		w := httptest.NewRecorder()

		mockUserUseCase.ExpectedCalls = nil
		mockUserUseCase.On("EditProfile",
			r.Context(),
			mock.AnythingOfType("models.User")).Return(item.mockUserUseCase...)

		userHandler.EditProfileHandler(w, r)

		CheckResponse(t, w, caseNum, item)
	}
}

func TestGetAllTags(t *testing.T) {
	t.Parallel()

	mockUserUseCase := &mocks.UserUsecase{}
	mockSessionUseCase := &_s.SessionUsecase{}

	call := mockUserUseCase.On("GetAllTags", context.Background())

	userHandler := &UserHandler{
		Logger:       logger.DripLogger,
		UserUCase:    mockUserUseCase,
		SessionUcase: mockSessionUseCase,
	}

	cases := []TestCase{
		{
			mockUserUseCase: []interface{}{
				tags,
				models.StatusOk200,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":{"allTags":` + tagsMapStr + `,"tagsCount":` + tagsCountStr + `}}`,
		},
		{
			mockUserUseCase: []interface{}{
				models.Tags{},
				models.HTTPError{
					Code: http.StatusNotFound,
				},
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
	}

	for caseNum, item := range cases {
		call.Return(item.mockUserUseCase...)

		r := httptest.NewRequest("GET", "/api/v1/tags", nil)
		w := httptest.NewRecorder()

		userHandler.GetAllTags(w, r)

		CheckResponse(t, w, caseNum, item)
	}
}

func TestMatches(t *testing.T) {
	t.Parallel()

	mockUserUseCase := &mocks.UserUsecase{}
	mockSessionUseCase := &_s.SessionUsecase{}

	call := mockUserUseCase.On("UsersMatches", context.Background())

	userHandler := &UserHandler{
		Logger:       logger.DripLogger,
		UserUCase:    mockUserUseCase,
		SessionUcase: mockSessionUseCase,
	}

	cases := []TestCase{
		{
			mockUserUseCase: []interface{}{
				matches,
				models.StatusOk200,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":{"allUsers":` + usersMapStr + `,"matchesCount":"` + matchesCountStr + `"}}`,
		},
		{
			mockUserUseCase: []interface{}{
				models.Matches{},
				models.HTTPError{
					Code: http.StatusNotFound,
				},
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
	}

	for caseNum, item := range cases {
		call.Return(item.mockUserUseCase...)

		r := httptest.NewRequest("GET", "/api/v1/matches", nil)
		w := httptest.NewRecorder()

		userHandler.MatchesHandler(w, r)

		CheckResponse(t, w, caseNum, item)
	}
}

func TestReaction(t *testing.T) {
	t.Parallel()

	mockUserUseCase := &mocks.UserUsecase{}
	mockSessionUseCase := &_s.SessionUsecase{}

	userHandler := &UserHandler{
		Logger:       logger.DripLogger,
		UserUCase:    mockUserUseCase,
		SessionUcase: mockSessionUseCase,
	}

	cases := []TestCase{
		{
			BodyReq: bytes.NewReader([]byte(`{"id":` + idStr + `,"reaction":` + reactionStr + `}`)),
			mockUserUseCase: []interface{}{
				match,
				models.StatusOk200,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":{"match":` + matchStr + `}}`,
		},
		{
			BodyReq: bytes.NewReader([]byte(`{"id":` + idStr + `,"reaction":` + reactionStr + `}`)),
			mockUserUseCase: []interface{}{
				notMatch,
				models.StatusOk200,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":{"match":` + notMatchStr + `}}`,
		},
		{
			BodyReq:    bytes.NewReader([]byte(`wrong input data`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":400,"body":null}`,
		},
		{
			BodyReq: bytes.NewReader([]byte(`{"id":` + idStr + `,"reaction":` + reactionStr + `}`)),
			mockUserUseCase: []interface{}{
				models.Match{},
				models.HTTPError{Code: http.StatusNotFound},
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
	}

	for caseNum, item := range cases {
		r := CreateRequest("POST", "/api/v1/likes", item.BodyReq)
		w := httptest.NewRecorder()

		mockUserUseCase.ExpectedCalls = nil
		mockUserUseCase.On("Reaction",
			r.Context(),
			mock.AnythingOfType("models.UserReaction")).Return(item.mockUserUseCase...)

		userHandler.ReactionHandler(w, r)

		CheckResponse(t, w, caseNum, item)
	}
}

func TestUploadPhoto(t *testing.T) {
	t.Parallel()

	mockUserUseCase := &mocks.UserUsecase{}
	mockSessionUseCase := &_s.SessionUsecase{}

	userHandler := &UserHandler{
		Logger:       logger.DripLogger,
		UserUCase:    mockUserUseCase,
		SessionUcase: mockSessionUseCase,
	}

	cases := []TestCase{
		{
			BodyReq: bytes.NewReader([]byte(`------boundary
Content-Disposition: form-data; name="photo"; filename="photo.jpg"
Content-Type: image/jpeg


------boundary--`)),
			mockUserUseCase: []interface{}{
				photo,
				models.StatusOk200,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":{"photo":"` + photo.Path + `"}}`,
		},
		{
			BodyReq:    bytes.NewReader([]byte(`wrong input data`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":400,"body":null}`,
		},
		{
			BodyReq: bytes.NewReader([]byte(`------boundary
Content-Disposition: form-data; name="wrong name"; filename="photo.jpg"
Content-Type: image/jpeg


------boundary--`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":400,"body":null}`,
		},
		{
			BodyReq: bytes.NewReader([]byte(`------boundary
Content-Disposition: form-data; name="photo"; filename="photo.jpg"
Content-Type: image/jpeg


------boundary--`)),
			mockUserUseCase: []interface{}{
				models.Photo{},
				models.HTTPError{Code: http.StatusInternalServerError},
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":500,"body":null}`,
		},
	}

	for caseNum, item := range cases {
		r := CreateRequest("POST", "/api/v1/profile/photo", item.BodyReq)
		r.Header.Add("Content-type", "multipart/form-data; boundary=----boundary")
		w := httptest.NewRecorder()

		mockUserUseCase.ExpectedCalls = nil
		mockUserUseCase.On("AddPhoto",
			r.Context(),
			mock.AnythingOfType("sectionReadCloser"),
			mock.AnythingOfType("string")).Return(item.mockUserUseCase...)

		userHandler.UploadPhoto(w, r)

		CheckResponse(t, w, caseNum, item)
	}
}

func TestDeletePhoto(t *testing.T) {
	t.Parallel()

	mockUserUseCase := &mocks.UserUsecase{}
	mockSessionUseCase := &_s.SessionUsecase{}

	userHandler := &UserHandler{
		Logger:       logger.DripLogger,
		UserUCase:    mockUserUseCase,
		SessionUcase: mockSessionUseCase,
	}

	cases := []TestCase{
		{
			BodyReq: bytes.NewReader([]byte(`{"photo":"` + photo.Path + `"}`)),
			mockUserUseCase: []interface{}{
				models.StatusOk200,
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":null}`,
		},
		{
			BodyReq:    bytes.NewReader([]byte(`wrong input data`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":400,"body":null}`,
		},
		{
			BodyReq: bytes.NewReader([]byte(`{"photo":"` + photo.Path + `"}`)),
			mockUserUseCase: []interface{}{
				models.HTTPError{Code: http.StatusBadRequest},
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":400,"body":null}`,
		},
	}

	for caseNum, item := range cases {
		r := CreateRequest("DELETE", "/api/v1/profile/photo", item.BodyReq)
		w := httptest.NewRecorder()

		mockUserUseCase.ExpectedCalls = nil
		mockUserUseCase.On("DeletePhoto",
			r.Context(),
			mock.AnythingOfType("models.Photo")).Return(item.mockUserUseCase...)

		userHandler.DeletePhoto(w, r)

		CheckResponse(t, w, caseNum, item)
	}
}

func TestSetRouting(t *testing.T) {
	mockUserUseCase := &mocks.UserUsecase{}
	mockSessionUseCase := &_s.SessionUsecase{}
	mockUserRepo := &repository.PostgreUserRepo{}

	SetUserRouting(logger.DripLogger, mux.NewRouter(), mockUserUseCase, mockSessionUseCase, mockUserRepo)
}
