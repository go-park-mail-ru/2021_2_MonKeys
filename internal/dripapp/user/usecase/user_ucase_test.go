package usecase_test

import (
	"context"
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	fileMocks "dripapp/internal/dripapp/file/mocks"
	userMocks "dripapp/internal/dripapp/user/mocks"
	"dripapp/internal/dripapp/user/usecase"

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

func TestUserUsecase_CurrentUser(t *testing.T) {
	type TestCase struct {
		userSession models.Session
		err         models.HTTPError
	}
	testCases := []TestCase{
		// Test OK
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			err: models.StatusOk200,
		},
		// Test ErrorNotFound
		{
			userSession: models.Session{
				UserID: 1,
				Cookie: "",
			},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "",
			},
		},
		// Test ErrContextNilError
		{
			userSession: models.Session{
				UserID: 2,
				Cookie: "",
			},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: models.ErrContextNilError,
			},
		},
	}

	type MockResultCase struct {
		user models.User
		err  error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			user: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			err: nil,
		},
		// Test ErrorNotFound
		{
			user: models.User{},
			err:  errors.New(""),
		},
		// Test ErrContextNilError
		{
			user: models.User{},
			err:  nil,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.userSession.UserID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ForContext, testCase.userSession))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetUserByID", r.Context(), mock.AnythingOfType("uint64")).Return(MockResultCases[i].user, MockResultCases[i].err)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2)
		user, status := testUserUsecase.CurrentUser(r.Context())

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(MockResultCases[i].user, user)

	}
}

func TestUserUsecase_EditProfile(t *testing.T) {
	type TestCase struct {
		userSession models.Session
		err         models.HTTPError
	}
	testCases := []TestCase{
		// Test OK
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			err: models.StatusOk200,
		},
		// Test ErrorNotFound
		{
			userSession: models.Session{
				UserID: 1,
				Cookie: "",
			},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "",
			},
		},
		// Test ErrContextNilError
		{
			userSession: models.Session{
				UserID: 2,
				Cookie: "",
			},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: models.ErrContextNilError,
			},
		},
		// Test ErrFailedToSaveAge
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			err: models.HTTPError{
				Code:    http.StatusBadRequest,
				Message: ("failed to save age"),
			},
		},
		// Test ErrUpdateUser
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			err: models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "",
			},
		},
		// Test ErrFailedToSaveAgeNewProfile
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			err: models.HTTPError{
				Code:    http.StatusBadRequest,
				Message: ("failed to save age"),
			},
		},
	}

	type MockResultCase struct {
		oldUser   models.User
		newUser   models.User
		errFirst  error
		errSecond error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			oldUser: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2001-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			newUser: models.User{
				ID:          0,
				Name:        "DripDrip",
				Email:       "drip@app.ru",
				Password:    "hahahahi",
				Date:        "2001-02-22",
				Age:         "19",
				Description: "vsem poka",
				Imgs:        []string{"1", "5"},
				Tags:        []string{"anime"},
			},
			errFirst:  nil,
			errSecond: nil,
		},
		// Test ErrorNotFound
		{
			oldUser:   models.User{},
			newUser:   models.User{},
			errFirst:  errors.New(""),
			errSecond: errors.New(""),
		},
		// Test ErrContextNilError
		{
			oldUser:   models.User{},
			newUser:   models.User{},
			errFirst:  nil,
			errSecond: nil,
		},
		// Test ErrFailedToSaveAge
		{
			oldUser: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "22-02-2001",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			newUser: models.User{
				ID:          0,
				Name:        "DripDrip",
				Email:       "drip@app.ru",
				Password:    "hahahahi",
				Date:        "22-02-2001",
				Age:         "19",
				Description: "vsem poka",
				Imgs:        []string{"1", "5"},
				Tags:        []string{"anime"},
			},
			errFirst:  nil,
			errSecond: nil,
		},
		// Test ErrUpdateUser
		{
			oldUser: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			newUser: models.User{
				ID:          0,
				Name:        "DripDrip",
				Email:       "drip@app.ru",
				Password:    "hahahahi",
				Date:        "2001-02-22",
				Age:         "19",
				Description: "vsem poka",
				Imgs:        []string{"1", "5"},
				Tags:        []string{"anime"},
			},
			errFirst:  nil,
			errSecond: errors.New(""),
		},
		// Test ErrFailedToSaveAgeNewProfile
		{
			oldUser: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			newUser: models.User{
				ID:          0,
				Name:        "DripDrip",
				Email:       "drip@app.ru",
				Password:    "hahahahi",
				Date:        "22-02-2001",
				Age:         "19",
				Description: "vsem poka",
				Imgs:        []string{"1", "5"},
				Tags:        []string{"anime"},
			},
			errFirst:  nil,
			errSecond: nil,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.userSession.UserID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ForContext, testCase.userSession))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetUserByID", r.Context(), mock.AnythingOfType("uint64")).Return(MockResultCases[i].oldUser, MockResultCases[i].errFirst)
		mockUserRepository.On("UpdateUser", r.Context(), mock.AnythingOfType("models.User")).Return(MockResultCases[i].newUser, MockResultCases[i].errSecond)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2)

		user, status := testUserUsecase.EditProfile(r.Context(), MockResultCases[i].newUser)

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(MockResultCases[i].newUser, user)

	}
}

// func TestUserUsecase_AddPhoto(t *testing.T) {
// 	type TestCase struct {
// 		userSession models.Session
// 		err         models.HTTPError
// 	}
// 	testCases := []TestCase{
// 		// Test OK
// 		{
// 			userSession: models.Session{
// 				UserID: 0,
// 				Cookie: "",
// 			},
// 			err: models.StatusOk200,
// 		},
// 		// Test ErrorNotFound
// 		{
// 			userSession: models.Session{
// 				UserID: 1,
// 				Cookie: "",
// 			},
// 			err: models.HTTPError{
// 				Code:    http.StatusNotFound,
// 				Message: "",
// 			},
// 		},
// 		// Test ErrContextNilError
// 		{
// 			userSession: models.Session{
// 				UserID: 2,
// 				Cookie: "",
// 			},
// 			err: models.HTTPError{
// 				Code:    http.StatusNotFound,
// 				Message: models.ErrContextNilError,
// 			},
// 		},
// 	}
//
// 	type MockResultCase struct {
// 		user models.User
// 		err  error
// 	}
// 	MockResultCases := []MockResultCase{
// 		// Test OK
// 		{
// 			user: models.User{
// 				ID:          0,
// 				Name:        "Drip",
// 				Email:       "drip@app.com",
// 				Password:    "hahaha",
// 				Date:        "2000-02-22",
// 				Description: "vsem privet",
// 				Imgs:        []string{"1", "2"},
// 				Tags:        []string{"anime", "BMSTU"},
// 			},
// 			err: nil,
// 		},
// 		// Test ErrorNotFound
// 		{
// 			user: models.User{},
// 			err:  errors.New(""),
// 		},
// 		// Test ErrContextNilError
// 		{
// 			user: models.User{},
// 			err:  nil,
// 		},
// 	}
//
// 	for i, testCase := range testCases {
// 		message := fmt.Sprintf("test case number: %d", i)
//
// 		// r, err := http.NewRequest(http.MethodGet, "test", nil)
// 		body := bytes.NewReader([]byte(`------boundary
// Content-Disposition: form-data; name="photo"; filename="photo.jpg"
// Content-Type: image/jpeg
//
// ------boundary--`))
// 		r, err := CreateMultipartRequest("POST", "/api/v1/profile/photo", body)
// 		r.Header.Add("Content-type", "multipart/form-data; boundary=----boundary")
// 		assert.NoError(t, err)
// 		if testCase.userSession.UserID != 2 {
// 			r = r.WithContext(context.WithValue(r.Context(), configs.ForContext, testCase.userSession))
// 		}
//
// 		mockUserRepository := new(userMocks.UserRepository)
// 		mockUserRepository.On("GetUserByID", r.Context(), mock.AnythingOfType("uint64")).Return(MockResultCases[i].user, MockResultCases[i].err)
// 		mockFileRepository := new(fileMocks.FileRepository)
//
// 		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2)
//
// 		uploadedPhoto, _, err := r.FormFile("photo")
// 		assert.NoError(t, err)
// 		defer uploadedPhoto.Close()
//
// 		user, status := testUserUsecase.AddPhoto(r.Context(), uploadedPhoto)
//
// 		assert.Equal(t, testCase.err, status, message)
// 		reflect.DeepEqual(MockResultCases[i].user, user)
//
// 	}
// }

func TestUserUsecase_Login(t *testing.T) {
	type TestCase struct {
		userSession models.Session
		logUserData models.LoginUser
		user        models.User
		err         models.HTTPError
	}
	testCases := []TestCase{
		// Test OK
		{
			userSession: models.Session{},
			logUserData: models.LoginUser{
				ID:       0,
				Email:    "drip@app.com",
				Password: "VBif222!",
			},
			user: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "oJnNPGsi805543a8fbee141b373962de3e347822de9ccb8e",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"}},
			err: models.StatusOk200,
		},
		// Test ErrorNotFoundEmail
		{
			userSession: models.Session{},
			logUserData: models.LoginUser{
				ID:       0,
				Email:    "drip@app.ru",
				Password: "VBif222!",
			},
			user: models.User{},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "",
			},
		},
		// Test ErrorNotFoundPassword
		{
			userSession: models.Session{},
			logUserData: models.LoginUser{
				ID:       0,
				Email:    "drip@app.com",
				Password: "VBif222!!",
			},
			user: models.User{},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "",
			},
		},
	}

	type MockResultCase struct {
		user models.User
		err  error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			user: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "oJnNPGsi805543a8fbee141b373962de3e347822de9ccb8e",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			err: nil,
		},
		// Test ErrorNotFoundEmail
		{
			user: models.User{},
			err:  errors.New(""),
		},
		// Test ErrorNotFoundPassword
		{
			user: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "oJnNPGsi805543a8fbee141b373962de3e347822de9ccb8e",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"}},
			err: nil,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		r = r.WithContext(context.WithValue(r.Context(), configs.ForContext, testCase.userSession))

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetUser", r.Context(), mock.AnythingOfType("string")).Return(MockResultCases[i].user, MockResultCases[i].err)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2)

		user, status := testUserUsecase.Login(r.Context(), testCase.logUserData)

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(testCase.user, user)

	}
}

func TestUserUsecase_Signup(t *testing.T) {
	type TestCase struct {
		userSession models.Session
		logUserData models.LoginUser
		user        models.User
		err         models.HTTPError
	}
	testCases := []TestCase{
		// Test OK
		{
			userSession: models.Session{},
			logUserData: models.LoginUser{
				ID:       0,
				Email:    "drip@app.com",
				Password: "VBif222!",
			},
			user: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "oJnNPGsi805543a8fbee141b373962de3e347822de9ccb8e",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			err: models.StatusOk200,
		},
		// Test ErrEmailAlreadyExists
		{
			userSession: models.Session{},
			logUserData: models.LoginUser{
				ID:       0,
				Email:    "drip@app.com",
				Password: "VBif222!",
			},
			user: models.User{},
			err: models.HTTPError{
				Code:    models.StatusEmailAlreadyExists,
				Message: "",
			},
		},
		// Test ErrCreateUser
		{
			userSession: models.Session{},
			logUserData: models.LoginUser{
				ID:       0,
				Email:    "drip@app.com",
				Password: "VBif222!",
			},
			user: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "oJnNPGsi805543a8fbee141b373962de3e347822de9ccb8e",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			err: models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "",
			},
		},
		// Test ErrCreateFolder
		{
			userSession: models.Session{},
			logUserData: models.LoginUser{
				ID:       0,
				Email:    "drip@app.com",
				Password: "VBif222!",
			},
			user: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "oJnNPGsi805543a8fbee141b373962de3e347822de9ccb8e",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			err: models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "",
			},
		},
	}

	type MockResultCase struct {
		curUser      models.User
		creatingUser models.User
		errDB        error
		errFile      error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			curUser: models.User{},
			creatingUser: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "oJnNPGsi805543a8fbee141b373962de3e347822de9ccb8e",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			errDB:   nil,
			errFile: nil,
		},
		// Test ErrEmailAlreadyExists
		{
			curUser: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "oJnNPGsi805543a8fbee141b373962de3e347822de9ccb8e",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			creatingUser: models.User{},
			errDB:        nil,
			errFile:      nil,
		},
		// Test ErrCreateUser
		{
			curUser: models.User{},
			creatingUser: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "oJnNPGsi805543a8fbee141b373962de3e347822de9ccb8e",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			errDB:   errors.New(""),
			errFile: nil,
		},
		// Test ErrCreateFolder
		{
			curUser: models.User{},
			creatingUser: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "oJnNPGsi805543a8fbee141b373962de3e347822de9ccb8e",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			errDB:   nil,
			errFile: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetUser", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("string")).Return(MockResultCases[i].curUser, MockResultCases[i].errDB)
		mockUserRepository.On("CreateUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("models.LoginUser")).Return(MockResultCases[i].creatingUser, MockResultCases[i].errDB)
		mockFileRepository := new(fileMocks.FileRepository)
		mockFileRepository.On("CreateFoldersForNewUser", mock.AnythingOfType("models.User")).Return(MockResultCases[i].errFile)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2)

		user, status := testUserUsecase.Signup(r.Context(), testCase.logUserData)

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(testCase.user, user)

	}
}

func TestUserUsecase_NextUser(t *testing.T) {
	type TestCase struct {
		userSession models.Session
		nextUsers   []models.User
		err         models.HTTPError
	}
	testCases := []TestCase{
		// Test OK
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			nextUsers: []models.User{
				{
					ID:          1,
					Name:        "Drip))",
					Email:       "drip@app.ru",
					Password:    "hah",
					Date:        "2001-02-22",
					Age:         "19",
					Description: "vsem privet",
					Imgs:        []string{"4", "3"},
					Tags:        []string{"BMSTU"},
				},
				{
					ID:          2,
					Name:        "Dr))",
					Email:       "dr@app.ru",
					Password:    "hah",
					Date:        "2000-02-22",
					Age:         "20",
					Description: "em privet",
					Imgs:        []string{"4", "3"},
					Tags:        []string{"JS"},
				},
				{
					ID:          3,
					Name:        "p))",
					Email:       "p@app.ru",
					Password:    "ah",
					Date:        "2001-02-22",
					Age:         "19",
					Description: "vs privet",
					Imgs:        []string{"3"},
					Tags:        []string{"BMSTU"},
				},
			},
			err: models.StatusOk200,
		},
		// Test ErrorNotFound
		{
			userSession: models.Session{
				UserID: 1,
				Cookie: "",
			},
			nextUsers: nil,
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "",
			},
		},
		// Test ErrContextNilError
		{
			userSession: models.Session{
				UserID: 2,
				Cookie: "",
			},
			nextUsers: nil,
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: models.ErrContextNilError,
			},
		},
		// Test ErrGetNextUserForSwipe
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			nextUsers: nil,
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: models.ErrContextNilError,
			},
		},
	}

	type MockResultCase struct {
		user            models.User
		nextUsers       []models.User
		errGetUser      error
		errGetNextUsers error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			user: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			nextUsers: []models.User{
				{
					ID:          1,
					Name:        "Drip))",
					Email:       "drip@app.ru",
					Password:    "hah",
					Date:        "2001-02-22",
					Age:         "19",
					Description: "vsem privet",
					Imgs:        []string{"4", "3"},
					Tags:        []string{"BMSTU"},
				},
				{
					ID:          2,
					Name:        "Dr))",
					Email:       "dr@app.ru",
					Password:    "hah",
					Date:        "2000-02-22",
					Age:         "20",
					Description: "em privet",
					Imgs:        []string{"4", "3"},
					Tags:        []string{"JS"},
				},
				{
					ID:          3,
					Name:        "p))",
					Email:       "p@app.ru",
					Password:    "ah",
					Date:        "2001-02-22",
					Age:         "19",
					Description: "vs privet",
					Imgs:        []string{"3"},
					Tags:        []string{"BMSTU"},
				},
			},
			errGetUser:      nil,
			errGetNextUsers: nil,
		},
		// Test ErrorNotFound
		{
			user:            models.User{},
			errGetUser:      errors.New(""),
			errGetNextUsers: nil,
		},
		// Test ErrContextNilError
		{
			user:            models.User{},
			errGetUser:      nil,
			errGetNextUsers: nil,
		},
		// Test ErrGetNextUserForSwipe
		{
			user: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			nextUsers:       nil,
			errGetUser:      nil,
			errGetNextUsers: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.userSession.UserID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ForContext, testCase.userSession))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetUserByID", r.Context(), mock.AnythingOfType("uint64")).Return(MockResultCases[i].user, MockResultCases[i].errGetUser)
		mockUserRepository.On("GetNextUserForSwipe", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("uint64")).Return(MockResultCases[i].nextUsers, MockResultCases[i].errGetNextUsers)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2)
		user, status := testUserUsecase.NextUser(r.Context())

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(testCase.nextUsers, user)

	}
}

func TestUserUsecase_GetAllTags(t *testing.T) {
	type TestCase struct {
		userSession models.Session
		tags        models.Tags
		err         models.HTTPError
	}
	testCases := []TestCase{
		// Test OK
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			tags: models.Tags{
				AllTags: map[uint64]models.Tag{
					0: models.Tag{Tag_Name: "anime"},
					1: models.Tag{Tag_Name: "BMSTU"},
					2: models.Tag{Tag_Name: "walk"},
					3: models.Tag{Tag_Name: "netflix"},
					4: models.Tag{Tag_Name: "prikolchiki"},
				},
				Count: 5,
			},
			err: models.StatusOk200,
		},
		// Test ErrorNotFound
		{
			userSession: models.Session{
				UserID: 1,
				Cookie: "",
			},
			tags: models.Tags{},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "",
			},
		},
		// Test ErrContextNilError
		{
			userSession: models.Session{
				UserID: 2,
				Cookie: "",
			},
			tags: models.Tags{},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: models.ErrContextNilError,
			},
		},
		// Test ErrGetTags
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			tags: models.Tags{},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "",
			},
		},
	}

	type MockResultCase struct {
		user       models.User
		tags       map[uint64]string
		errGetUser error
		errGetTags error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			user: models.User{},
			tags: map[uint64]string{
				0: "anime",
				1: "BMSTU",
				2: "walk",
				3: "netflix",
				4: "prikolchiki",
			},
			errGetUser: nil,
			errGetTags: nil,
		},
		// Test ErrorNotFound
		{
			user:       models.User{},
			tags:       map[uint64]string{},
			errGetUser: errors.New(""),
			errGetTags: nil,
		},
		// Test ErrContextNilError
		{
			user:       models.User{},
			tags:       map[uint64]string{},
			errGetUser: nil,
			errGetTags: nil,
		},
		// Test ErrGetNextUserForSwipe
		{
			user: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			tags:       map[uint64]string{},
			errGetUser: nil,
			errGetTags: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.userSession.UserID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ForContext, testCase.userSession))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetUserByID", r.Context(), mock.AnythingOfType("uint64")).Return(MockResultCases[i].user, MockResultCases[i].errGetUser)
		mockUserRepository.On("GetTags", mock.AnythingOfType("*context.timerCtx")).Return(MockResultCases[i].tags, MockResultCases[i].errGetTags)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2)
		allTags, status := testUserUsecase.GetAllTags(r.Context())

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(testCase.tags, allTags)

	}
}

func TestUserUsecase_UsersMatches(t *testing.T) {
	type TestCase struct {
		userSession models.Session
		matches     models.Matches
		err         models.HTTPError
	}
	testCases := []TestCase{
		// Test OK
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			matches: models.Matches{
				AllUsers: map[uint64]models.User{
					0: {
						ID:          1,
						Name:        "Drip))",
						Email:       "drip@app.ru",
						Password:    "hah",
						Date:        "2001-02-22",
						Age:         "19",
						Description: "vsem privet",
						Imgs:        []string{"4", "3"},
						Tags:        []string{"BMSTU"},
					},
					1: {
						ID:          2,
						Name:        "Dr))",
						Email:       "dr@app.ru",
						Password:    "hah",
						Date:        "2000-02-22",
						Age:         "20",
						Description: "em privet",
						Imgs:        []string{"4", "3"},
						Tags:        []string{"JS"},
					},
				},
				Count: "2",
			},
			err: models.StatusOk200,
		},
		// Test ErrorNotFound
		{
			userSession: models.Session{
				UserID: 1,
				Cookie: "",
			},
			matches: models.Matches{},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "",
			},
		},
		// Test ErrContextNilError
		{
			userSession: models.Session{
				UserID: 2,
				Cookie: "",
			},
			matches: models.Matches{},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: models.ErrContextNilError,
			},
		},
		// Test ErrGetTags
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			matches: models.Matches{},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "",
			},
		},
	}

	type MockResultCase struct {
		user          models.User
		matches       []models.User
		errGetUser    error
		errGetMatches error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			user: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			matches: []models.User{
				{
					ID:          1,
					Name:        "Drip))",
					Email:       "drip@app.ru",
					Password:    "hah",
					Date:        "2001-02-22",
					Age:         "19",
					Description: "vsem privet",
					Imgs:        []string{"4", "3"},
					Tags:        []string{"BMSTU"},
				},
				{
					ID:          2,
					Name:        "Dr))",
					Email:       "dr@app.ru",
					Password:    "hah",
					Date:        "2000-02-22",
					Age:         "20",
					Description: "em privet",
					Imgs:        []string{"4", "3"},
					Tags:        []string{"JS"},
				},
			},
			errGetUser:    nil,
			errGetMatches: nil,
		},
		// Test ErrorNotFound
		{
			user:          models.User{},
			matches:       []models.User{},
			errGetUser:    errors.New(""),
			errGetMatches: nil,
		},
		// Test ErrContextNilError
		{
			user:          models.User{},
			matches:       []models.User{},
			errGetUser:    nil,
			errGetMatches: nil,
		},
		// Test ErrGetNextUserForSwipe
		{
			user: models.User{
				ID:          0,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			matches:       []models.User{},
			errGetUser:    nil,
			errGetMatches: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.userSession.UserID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ForContext, testCase.userSession))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetUserByID", r.Context(), mock.AnythingOfType("uint64")).Return(MockResultCases[i].user, MockResultCases[i].errGetUser)
		mockUserRepository.On("GetUsersMatches", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("uint64")).Return(MockResultCases[i].matches, MockResultCases[i].errGetMatches)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2)
		allMatches, status := testUserUsecase.UsersMatches(r.Context())

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(testCase.matches, allMatches)

	}
}

func TestUserUsecase_Reaction(t *testing.T) {
	type TestCase struct {
		userSession  models.Session
		reactionData models.UserReaction
		match        models.Match
		err          models.HTTPError
	}
	testCases := []TestCase{
		// Test OK and Match
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			reactionData: models.UserReaction{
				Id:       2,
				Reaction: 1,
			},
			match: models.Match{
				Match: true,
			},
			err: models.StatusOk200,
		},
		// Test OK and no Match
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			reactionData: models.UserReaction{
				Id:       5,
				Reaction: 1,
			},
			match: models.Match{
				Match: false,
			},
			err: models.StatusOk200,
		},
		// Test OK and no Match
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			reactionData: models.UserReaction{
				Id:       2,
				Reaction: 2,
			},
			match: models.Match{
				Match: false,
			},
			err: models.StatusOk200,
		},
		// Test ErrorNotFound
		{
			userSession: models.Session{
				UserID: 1,
				Cookie: "",
			},
			reactionData: models.UserReaction{},
			match:        models.Match{},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "",
			},
		},
		// Test ErrContextNilError
		{
			userSession: models.Session{
				UserID: 2,
				Cookie: "",
			},
			reactionData: models.UserReaction{},
			match:        models.Match{},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: models.ErrContextNilError,
			},
		},
		// Test ErrAddReaction
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			reactionData: models.UserReaction{
				Id:       2,
				Reaction: 1,
			},
			match: models.Match{},
			err: models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "",
			},
		},
		// Test ErrGetLikes
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			reactionData: models.UserReaction{
				Id:       2,
				Reaction: 1,
			},
			match: models.Match{},
			err: models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "",
			},
		},
		// Test ErrDeleteLike
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			reactionData: models.UserReaction{
				Id:       2,
				Reaction: 1,
			},
			match: models.Match{},
			err: models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "",
			},
		},
		// Test ErrAddMatch
		{
			userSession: models.Session{
				UserID: 0,
				Cookie: "",
			},
			reactionData: models.UserReaction{
				Id:       2,
				Reaction: 1,
			},
			match: models.Match{},
			err: models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "",
			},
		},
	}

	type MockResultCase struct {
		user           models.User
		likes          []uint64
		errGetUser     error
		errAddReaction error
		errGetLikes    error
		errDeleteLike  error
		errAddMatch    error
	}
	MockResultCases := []MockResultCase{
		// Test OK and Match
		{
			user:           models.User{},
			likes:          []uint64{1, 2, 3},
			errGetUser:     nil,
			errAddReaction: nil,
			errGetLikes:    nil,
			errDeleteLike:  nil,
			errAddMatch:    nil,
		},
		// Test OK and no Match
		{
			user:           models.User{},
			likes:          []uint64{1, 2, 3},
			errGetUser:     nil,
			errAddReaction: nil,
			errGetLikes:    nil,
			errDeleteLike:  nil,
			errAddMatch:    nil,
		},
		// Test OK and no Match
		{
			user:           models.User{},
			likes:          []uint64{1, 2, 3},
			errGetUser:     nil,
			errAddReaction: nil,
			errGetLikes:    nil,
			errDeleteLike:  nil,
			errAddMatch:    nil,
		},
		// Test ErrorNotFound
		{
			user:           models.User{},
			likes:          []uint64{},
			errGetUser:     errors.New(""),
			errAddReaction: nil,
			errGetLikes:    nil,
			errDeleteLike:  nil,
			errAddMatch:    nil,
		},
		// Test ErrContextNilError
		{
			user:           models.User{},
			likes:          []uint64{},
			errGetUser:     nil,
			errAddReaction: nil,
			errGetLikes:    nil,
			errDeleteLike:  nil,
			errAddMatch:    nil,
		},
		// Test ErrAddReaction
		{
			user:           models.User{},
			likes:          []uint64{},
			errGetUser:     nil,
			errAddReaction: errors.New(""),
			errGetLikes:    nil,
			errDeleteLike:  nil,
			errAddMatch:    nil,
		},
		// Test ErrGetLikes
		{
			user:           models.User{},
			likes:          []uint64{},
			errGetUser:     nil,
			errAddReaction: nil,
			errGetLikes:    errors.New(""),
			errDeleteLike:  nil,
			errAddMatch:    nil,
		},
		// Test ErrDeleteLike
		{
			user:           models.User{},
			likes:          []uint64{1, 2, 3},
			errGetUser:     nil,
			errAddReaction: nil,
			errGetLikes:    nil,
			errDeleteLike:  errors.New(""),
			errAddMatch:    nil,
		},
		// Test ErrAddMatch
		{
			user:           models.User{},
			likes:          []uint64{1, 2, 3},
			errGetUser:     nil,
			errAddReaction: nil,
			errGetLikes:    nil,
			errDeleteLike:  nil,
			errAddMatch:    errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.userSession.UserID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ForContext, testCase.userSession))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetUserByID",
			r.Context(),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].user, MockResultCases[i].errGetUser)
		mockUserRepository.On("AddReaction",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].errAddReaction)
		mockUserRepository.On("GetLikes",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].likes, MockResultCases[i].errGetLikes)
		mockUserRepository.On("DeleteLike",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].errDeleteLike)
		mockUserRepository.On("AddMatch",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].errAddMatch)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2)
		match, status := testUserUsecase.Reaction(r.Context(), testCase.reactionData)

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(testCase.match, match)

	}
}
