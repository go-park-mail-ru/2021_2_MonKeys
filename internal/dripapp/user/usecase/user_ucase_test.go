package usecase_test

import (
	"bytes"
	"context"
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	"errors"
	"fmt"
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

func TestUserUsecase_CurrentUser(t *testing.T) {
	type TestCase struct {
		user models.User
		err  error
	}
	testCases := []TestCase{
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
		// // Test ErrContextNilError
		{
			user: models.User{
				ID: 2,
			},
			err: models.ErrContextNilError,
		},
	}

	type MockResultCase struct {
		user models.User
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
		},
		// Test ErrContextNilError
		{
			user: models.User{
				ID: 2,
			},
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.user.ID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, testCase.user))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2, nil)
		user, status := testUserUsecase.CurrentUser(r.Context())

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(MockResultCases[i].user, user)

	}
}

func TestUserUsecase_EditProfile(t *testing.T) {
	type TestCase struct {
		user models.User
		err  error
	}
	testCases := []TestCase{
		// Test OK
		{
			user: models.User{
				ID: 0,
			},
			err: nil,
		},
		// Test ErrorNotFound
		{
			user: models.User{
				ID: 1,
			},
			err: errors.New("failed on userYear"),
		},
		// Test ErrContextNilError
		{
			user: models.User{
				ID: 2,
			},
			err: models.ErrContextNilError,
		},
		// Test ErrUpdateUser
		{
			user: models.User{
				ID: 0,
			},
			err: errors.New(""),
		},
	}

	type MockResultCase struct {
		newUser models.User
		err     error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			newUser: models.User{
				ID:          1,
				Name:        "DripDrip",
				Email:       "drip@app.ru",
				Password:    "hahahahi",
				Date:        "2001-02-22",
				Age:         19,
				Description: "vsem poka",
				Imgs:        []string{"1", "5"},
				Tags:        []string{"anime"},
			},
			err: nil,
		},
		// Test ErrorNotFound
		{
			newUser: models.User{},
			err:     errors.New("failed on userYear"),
		},
		// Test ErrContextNilError
		{
			newUser: models.User{},
			err:     nil,
		},
		// Test ErrUpdateUser
		{
			newUser: models.User{
				ID:          0,
				Name:        "DripDrip",
				Email:       "drip@app.ru",
				Password:    "hahahahi",
				Date:        "2001-02-22",
				Age:         19,
				Description: "vsem poka",
				Imgs:        []string{"1", "5"},
				Tags:        []string{"anime"},
			},
			err: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.user.ID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, testCase.user))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("UpdateUser",
			r.Context(),
			mock.AnythingOfType("models.User")).Return(MockResultCases[i].newUser, MockResultCases[i].err)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2, nil)

		user, status := testUserUsecase.EditProfile(r.Context(), MockResultCases[i].newUser)

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(MockResultCases[i].newUser, user)

	}
}

func TestUserUsecase_AddPhoto(t *testing.T) {
	type TestCase struct {
		user models.User
		path string
		err  error
	}
	testCases := []TestCase{
		// Test OK
		{
			user: models.User{
				ID: 1,
			},
			path: "",
			err:  nil,
		},
		// Test ErrorNotFound
		{
			user: models.User{
				ID: 1,
			},
			err: nil,
		},
		// Test ErrContextNilError
		{
			user: models.User{
				ID: 2,
			},
			err: models.ErrContextNilError,
		},
		// Test ErrSaveUserPhoto
		{
			user: models.User{
				ID: 0,
			},
			err: errors.New(""),
		},
		// Test ErrUpdateImgs
		{
			user: models.User{
				ID: 0,
			},
			err: errors.New(""),
		},
	}

	type MockResultCase struct {
		path          string
		errSavePhoto  error
		errUpdateImgs error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			path:          "",
			errSavePhoto:  nil,
			errUpdateImgs: nil,
		},
		// Test ErrorNotFound
		{
			errSavePhoto:  nil,
			errUpdateImgs: nil,
		},
		// Test ErrContextNilError
		{
			errSavePhoto:  nil,
			errUpdateImgs: nil,
		},
		// Test ErrSaveUserPhoto
		{
			errSavePhoto:  errors.New(""),
			errUpdateImgs: nil,
		},
		// Test ErrUpdateImgs
		{
			errSavePhoto:  nil,
			errUpdateImgs: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		body := bytes.NewReader([]byte(`------boundary
Content-Disposition: form-data; name="photo"; filename="photo.jpg"
Content-Type: image/jpeg

------boundary--`))
		r, err := http.NewRequest("POST", "test", body)
		assert.NoError(t, err)
		r.Header.Add("Content-type", "multipart/form-data; boundary=----boundary")
		if testCase.user.ID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, testCase.user))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("UpdateImgs",
			r.Context(),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("[]string")).Return(MockResultCases[i].errUpdateImgs)
		mockFileRepository := new(fileMocks.FileRepository)
		mockFileRepository.On("SaveUserPhoto",
			mock.AnythingOfType("models.User"),
			mock.AnythingOfType("multipart.sectionReadCloser"),
			mock.AnythingOfType("string")).Return(MockResultCases[i].path, MockResultCases[i].errSavePhoto)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2, nil)

		uploadedPhoto, _, err := r.FormFile("photo")
		assert.NoError(t, err)
		defer uploadedPhoto.Close()

		path, status := testUserUsecase.AddPhoto(r.Context(), uploadedPhoto, testCase.path)

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(MockResultCases[i].path, path)

	}
}

func TestUserUsecase_DeletePhoto(t *testing.T) {
	type TestCase struct {
		user  models.User
		photo models.Photo
		err   error
	}
	testCases := []TestCase{
		// Test OK
		{
			user: models.User{
				ID:          1,
				Name:        "Drip",
				Email:       "drip@app.com",
				Password:    "hahaha",
				Date:        "2000-02-22",
				Description: "vsem privet",
				Imgs:        []string{"1", "2"},
				Tags:        []string{"anime", "BMSTU"},
			},
			photo: models.Photo{
				Path: "1",
			},
			err: nil,
		},
		// Test ErrContextNilError
		{
			user: models.User{
				ID: 2,
			},
			photo: models.Photo{},
			err:   models.ErrContextNilError,
		},
		// Test ErrDelete
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
			photo: models.Photo{
				Path: "",
			},
			err: errors.New(""),
		},
		// Test ErrUpdateImgs
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
			photo: models.Photo{
				Path: "",
			},
			err: errors.New(""),
		},
		// Test ErrDeletePhoto
		{
			user: models.User{},
			photo: models.Photo{
				Path: "",
			},
			err: errors.New("user does not have such a photo"),
		},
	}

	type MockResultCase struct {
		errDelete     error
		errUpdateImgs error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			errDelete:     nil,
			errUpdateImgs: nil,
		},
		// Test ErrContextNilError
		{
			errDelete:     nil,
			errUpdateImgs: nil,
		},
		// Test ErrDelete
		{
			errDelete:     errors.New(""),
			errUpdateImgs: nil,
		},
		// Test ErrUpdateImgs
		{
			errDelete:     nil,
			errUpdateImgs: errors.New(""),
		},
		// Test ErrDeletePhoto
		{
			errDelete:     nil,
			errUpdateImgs: nil,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		body := bytes.NewReader([]byte(`------boundary
Content-Disposition: form-data; name="photo"; filename="photo.jpg"
Content-Type: image/jpeg

------boundary--`))
		r, err := http.NewRequest("POST", "test", body)
		assert.NoError(t, err)
		r.Header.Add("Content-type", "multipart/form-data; boundary=----boundary")
		if testCase.user.ID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, testCase.user))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("UpdateImgs",
			r.Context(),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("[]string")).Return(MockResultCases[i].errUpdateImgs)
		mockFileRepository := new(fileMocks.FileRepository)
		mockFileRepository.On("Delete",
			mock.AnythingOfType("string")).Return(MockResultCases[i].errDelete)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2, nil)

		status := testUserUsecase.DeletePhoto(r.Context(), testCase.photo)

		assert.Equal(t, testCase.err, status, message)
	}
}

func TestUserUsecase_Login(t *testing.T) {
	type TestCase struct {
		logUserData models.LoginUser
		user        models.User
		err         error
	}
	testCases := []TestCase{
		// Test OK
		{
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
			err: nil,
		},
		// Test ErrorNotFoundEmail
		{
			logUserData: models.LoginUser{
				ID:       0,
				Email:    "drip@app.ru",
				Password: "VBif222!",
			},
			user: models.User{},
			err:  errors.New(""),
		},
		// Test ErrorNotFoundPassword
		{
			logUserData: models.LoginUser{
				ID:       0,
				Email:    "drip@app.com",
				Password: "VBif222!!",
			},
			user: models.User{},
			err:  models.ErrMismatch,
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

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetUser",
			r.Context(),
			mock.AnythingOfType("string")).Return(MockResultCases[i].user, MockResultCases[i].err)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2, nil)

		user, status := testUserUsecase.Login(r.Context(), testCase.logUserData)

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(testCase.user, user)

	}
}
func TestUserUsecase_Signup(t *testing.T) {
	type TestCase struct {
		logUserData models.LoginUser
		user        models.User
		err         error
	}
	testCases := []TestCase{
		// Test OK
		{
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
			err: nil,
		},
		// Test ErrEmailAlreadyExists
		{
			logUserData: models.LoginUser{
				ID:       0,
				Email:    "drip@app.com",
				Password: "VBif222!",
			},
			user: models.User{},
			err:  models.ErrEmailAlreadyExists,
		},
		// Test ErrCreateUser
		{
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
			err: errors.New(""),
		},
		// Test ErrCreateFolder
		{
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
			err: errors.New(""),
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
		mockUserRepository.On("GetUser",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("string")).Return(MockResultCases[i].curUser, MockResultCases[i].errDB)
		mockUserRepository.On("CreateUser",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("models.LoginUser")).Return(MockResultCases[i].creatingUser, MockResultCases[i].errDB)
		mockFileRepository := new(fileMocks.FileRepository)
		mockFileRepository.On("CreateFoldersForNewUser",
			mock.AnythingOfType("models.User")).Return(MockResultCases[i].errFile)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2, nil)

		user, err := testUserUsecase.Signup(r.Context(), testCase.logUserData)

		assert.Equal(t, testCase.err, err, message)
		reflect.DeepEqual(testCase.user, user)
	}
}

func TestUserUsecase_NextUser(t *testing.T) {
	type TestCase struct {
		user      models.User
		nextUsers []models.User
		err       error
	}
	testCases := []TestCase{
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
					Age:         19,
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
					Age:         20,
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
					Age:         19,
					Description: "vs privet",
					Imgs:        []string{"3"},
					Tags:        []string{"BMSTU"},
				},
			},
			err: nil,
		},
		// Test ErrorNotFound
		{
			user:      models.User{},
			nextUsers: nil,
			err:       nil,
		},
		// Test ErrContextNilError
		{
			user: models.User{
				ID: 2,
			},
			nextUsers: nil,
			err:       models.ErrContextNilError,
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
			nextUsers: nil,
			err:       errors.New(""),
		},
	}

	type MockResultCase struct {
		nextUsers       []models.User
		errGetNextUsers error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			nextUsers: []models.User{
				{
					ID:          1,
					Name:        "Drip))",
					Email:       "drip@app.ru",
					Password:    "hah",
					Date:        "2001-02-22",
					Age:         19,
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
					Age:         20,
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
					Age:         19,
					Description: "vs privet",
					Imgs:        []string{"3"},
					Tags:        []string{"BMSTU"},
				},
			},
			errGetNextUsers: nil,
		},
		// Test ErrorNotFound
		{
			errGetNextUsers: nil,
		},
		// Test ErrContextNilError
		{
			errGetNextUsers: nil,
		},
		// Test ErrGetNextUserForSwipe
		{
			nextUsers:       nil,
			errGetNextUsers: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.user.ID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, testCase.user))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetNextUserForSwipe",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("models.User")).Return(MockResultCases[i].nextUsers, MockResultCases[i].errGetNextUsers)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2, nil)
		user, status := testUserUsecase.NextUser(r.Context())

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(testCase.nextUsers, user)

	}
}

func TestUserUsecase_GetAllTags(t *testing.T) {
	type TestCase struct {
		user models.User
		tags models.Tags
		err  error
	}
	testCases := []TestCase{
		// Test OK
		{
			user: models.User{},
			tags: models.Tags{
				AllTags: map[uint64]models.Tag{
					0: {TagName: "anime"},
					1: {TagName: "BMSTU"},
					2: {TagName: "walk"},
					3: {TagName: "netflix"},
					4: {TagName: "prikolchiki"},
				},
				Count: 5,
			},
			err: nil,
		},
		// Test ErrGetTags
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
			tags: models.Tags{},
			err:  errors.New(""),
		},
	}

	type MockResultCase struct {
		tags       map[uint64]string
		errGetTags error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			tags: map[uint64]string{
				0: "anime",
				1: "BMSTU",
				2: "walk",
				3: "netflix",
				4: "prikolchiki",
			},
			errGetTags: nil,
		},
		// Test ErrGetTags
		{
			tags:       map[uint64]string{},
			errGetTags: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.user.ID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, testCase.user))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetTags",
			mock.AnythingOfType("*context.timerCtx")).Return(MockResultCases[i].tags, MockResultCases[i].errGetTags)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2, nil)
		allTags, status := testUserUsecase.GetAllTags(r.Context())

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(testCase.tags, allTags)

	}
}

func TestUserUsecase_UsersMatches(t *testing.T) {
	type TestCase struct {
		user    models.User
		matches models.Matches
		err     error
	}
	testCases := []TestCase{
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
			matches: models.Matches{
				AllUsers: map[uint64]models.User{
					0: {
						ID:          1,
						Name:        "Drip))",
						Email:       "drip@app.ru",
						Password:    "hah",
						Date:        "2001-02-22",
						Age:         19,
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
						Age:         20,
						Description: "em privet",
						Imgs:        []string{"4", "3"},
						Tags:        []string{"JS"},
					},
				},
				Count: "2",
			},
			err: nil,
		},
		// Test ErrorNotFound
		{
			user:    models.User{},
			matches: models.Matches{},
			err:     nil,
		},
		// Test ErrContextNilError
		{
			user: models.User{
				ID: 2,
			},
			matches: models.Matches{},
			err:     models.ErrContextNilError,
		},
		// Test ErrGetTags
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
			matches: models.Matches{},
			err:     errors.New(""),
		},
	}

	type MockResultCase struct {
		matches       []models.User
		errGetMatches error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			matches: []models.User{
				{
					ID:          1,
					Name:        "Drip))",
					Email:       "drip@app.ru",
					Password:    "hah",
					Date:        "2001-02-22",
					Age:         19,
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
					Age:         20,
					Description: "em privet",
					Imgs:        []string{"4", "3"},
					Tags:        []string{"JS"},
				},
			},
			errGetMatches: nil,
		},
		// Test ErrorNotFound
		{
			matches:       []models.User{},
			errGetMatches: nil,
		},
		// Test ErrContextNilError
		{
			matches:       []models.User{},
			errGetMatches: nil,
		},
		// Test ErrGetNextUserForSwipe
		{
			matches:       []models.User{},
			errGetMatches: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.user.ID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, testCase.user))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetUsersMatches",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].matches, MockResultCases[i].errGetMatches)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2, nil)
		allMatches, status := testUserUsecase.UsersMatches(r.Context())

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(testCase.matches, allMatches)

	}
}

// func TestUserUsecase_Reaction(t *testing.T) {
// 	type TestCase struct {
// 		user         models.User
// 		reactionData models.UserReaction
// 		match        models.Match
// 		err          error
// 	}
// 	testCases := []TestCase{
// 		// Test OK and Match
// 		{
// 			user: models.User{
// 				ID: 0,
// 			},
// 			reactionData: models.UserReaction{
// 				Id:       2,
// 				Reaction: 1,
// 			},
// 			match: models.Match{
// 				Match: true,
// 			},
// 			err: nil,
// 		},
// 		// Test OK and no Match
// 		{
// 			user: models.User{
// 				ID: 0,
// 			},
// 			reactionData: models.UserReaction{
// 				Id:       5,
// 				Reaction: 1,
// 			},
// 			match: models.Match{
// 				Match: false,
// 			},
// 			err: nil,
// 		},
// 		// Test OK and no Match
// 		{
// 			user: models.User{
// 				ID: 0,
// 			},
// 			reactionData: models.UserReaction{
// 				Id:       2,
// 				Reaction: 2,
// 			},
// 			match: models.Match{
// 				Match: false,
// 			},
// 			err: nil,
// 		},
// 		// Test ErrorNotFound
// 		{
// 			user: models.User{
// 				ID: 1,
// 			},
// 			reactionData: models.UserReaction{},
// 			match:        models.Match{},
// 			err:          nil,
// 		},
// 		// Test ErrContextNilError
// 		{
// 			user: models.User{
// 				ID: 2,
// 			},
// 			reactionData: models.UserReaction{},
// 			match:        models.Match{},
// 			err:          models.ErrContextNilError,
// 		},
// 		// Test ErrAddReaction
// 		{
// 			user: models.User{
// 				ID: 0,
// 			},
// 			reactionData: models.UserReaction{
// 				Id:       2,
// 				Reaction: 1,
// 			},
// 			match: models.Match{},
// 			err:   errors.New(""),
// 		},
// 		// Test ErrGetLikes
// 		{
// 			user: models.User{
// 				ID: 0,
// 			},
// 			reactionData: models.UserReaction{
// 				Id:       2,
// 				Reaction: 1,
// 			},
// 			match: models.Match{},
// 			err:   errors.New(""),
// 		},
// 		// Test ErrDeleteLike
// 		{
// 			user: models.User{
// 				ID: 0,
// 			},
// 			reactionData: models.UserReaction{
// 				Id:       2,
// 				Reaction: 1,
// 			},
// 			match: models.Match{},
// 			err:   errors.New(""),
// 		},
// 		// Test ErrAddMatch
// 		{
// 			user: models.User{
// 				ID: 0,
// 			},
// 			reactionData: models.UserReaction{
// 				Id:       2,
// 				Reaction: 1,
// 			},
// 			match: models.Match{},
// 			err:   errors.New(""),
// 		},
// 	}

// 	type MockResultCase struct {
// 		likes             []uint64
// 		errAddReaction    error
// 		errGetLikes       error
// 		errDeleteReaction error
// 		errAddMatch       error
// 	}
// 	MockResultCases := []MockResultCase{
// 		// Test OK and Match
// 		{
// 			likes:             []uint64{1, 2, 3},
// 			errAddReaction:    nil,
// 			errGetLikes:       nil,
// 			errDeleteReaction: nil,
// 			errAddMatch:       nil,
// 		},
// 		// Test OK and no Match
// 		{
// 			likes:             []uint64{1, 2, 3},
// 			errAddReaction:    nil,
// 			errGetLikes:       nil,
// 			errDeleteReaction: nil,
// 			errAddMatch:       nil,
// 		},
// 		// Test OK and no Match
// 		{
// 			likes:             []uint64{1, 2, 3},
// 			errAddReaction:    nil,
// 			errGetLikes:       nil,
// 			errDeleteReaction: nil,
// 			errAddMatch:       nil,
// 		},
// 		// Test ErrorNotFound
// 		{
// 			likes:             []uint64{},
// 			errAddReaction:    nil,
// 			errGetLikes:       nil,
// 			errDeleteReaction: nil,
// 			errAddMatch:       nil,
// 		},
// 		// Test ErrContextNilError
// 		{
// 			likes:             []uint64{},
// 			errAddReaction:    nil,
// 			errGetLikes:       nil,
// 			errDeleteReaction: nil,
// 			errAddMatch:       nil,
// 		},
// 		// Test ErrAddReaction
// 		{
// 			likes:             []uint64{},
// 			errAddReaction:    errors.New(""),
// 			errGetLikes:       nil,
// 			errDeleteReaction: nil,
// 			errAddMatch:       nil,
// 		},
// 		// Test ErrGetLikes
// 		{
// 			likes:             []uint64{},
// 			errAddReaction:    nil,
// 			errGetLikes:       errors.New(""),
// 			errDeleteReaction: nil,
// 			errAddMatch:       nil,
// 		},
// 		// Test ErrDeleteLike
// 		{
// 			likes:             []uint64{1, 2, 3},
// 			errAddReaction:    nil,
// 			errGetLikes:       nil,
// 			errDeleteReaction: errors.New(""),
// 			errAddMatch:       nil,
// 		},
// 		// Test ErrAddMatch
// 		{
// 			likes:             []uint64{1, 2, 3},
// 			errAddReaction:    nil,
// 			errGetLikes:       nil,
// 			errDeleteReaction: nil,
// 			errAddMatch:       errors.New(""),
// 		},
// 	}

// 	for i, testCase := range testCases {
// 		message := fmt.Sprintf("test case number: %d", i)

// 		r, err := http.NewRequest(http.MethodGet, "test", nil)
// 		assert.NoError(t, err)
// 		if testCase.user.ID != 2 {
// 			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, testCase.user))
// 		}

// 		mockUserRepository := new(userMocks.UserRepository)
// 		mockUserRepository.On("AddReaction",
// 			mock.AnythingOfType("*context.timerCtx"),
// 			mock.AnythingOfType("uint64"),
// 			mock.AnythingOfType("uint64"),
// 			mock.AnythingOfType("uint64")).Return(MockResultCases[i].errAddReaction)
// 		mockUserRepository.On("GetLikes",
// 			mock.AnythingOfType("*context.timerCtx"),
// 			mock.AnythingOfType("uint64")).Return(MockResultCases[i].likes, MockResultCases[i].errGetLikes)
// 		mockUserRepository.On("DeleteReaction",
// 			mock.AnythingOfType("*context.timerCtx"),
// 			mock.AnythingOfType("uint64"),
// 			mock.AnythingOfType("uint64")).Return(MockResultCases[i].errDeleteReaction)
// 		mockUserRepository.On("AddMatch",
// 			mock.AnythingOfType("*context.timerCtx"),
// 			mock.AnythingOfType("uint64"),
// 			mock.AnythingOfType("uint64")).Return(MockResultCases[i].errAddMatch)
// 		mockFileRepository := new(fileMocks.FileRepository)

// 		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2)
// 		match, status := testUserUsecase.Reaction(r.Context(), testCase.reactionData)

// 		assert.Equal(t, testCase.err, status, message)
// 		reflect.DeepEqual(testCase.match, match)

// 	}
// }

func TestUserUsecase_UsersMatchesWithSearching(t *testing.T) {
	type TestCase struct {
		user          models.User
		matches       models.Matches
		searchingData models.Search
		err           error
	}
	testCases := []TestCase{
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
			matches: models.Matches{
				AllUsers: map[uint64]models.User{
					0: {
						ID:          1,
						Name:        "Владимир",
						Email:       "drip@app.ru",
						Password:    "hah",
						Date:        "2001-02-22",
						Age:         19,
						Description: "vsem privet",
						Imgs:        []string{"4", "3"},
						Tags:        []string{"BMSTU"},
					},
					1: {
						ID:          2,
						Name:        "влаД",
						Email:       "dr@app.ru",
						Password:    "hah",
						Date:        "2000-02-22",
						Age:         20,
						Description: "em privet",
						Imgs:        []string{"4", "3"},
						Tags:        []string{"JS"},
					},
				},
				Count: "2",
			},
			searchingData: models.Search{SearchingTmpl: "влад"},
			err:           nil,
		},
		// Test ErrContextNilError
		{
			user: models.User{
				ID: 2,
			},
			err: models.ErrContextNilError,
		},
		// Test ErrorNotFound
		{
			user:    models.User{},
			matches: models.Matches{},
			err:     errors.New(""),
		},
	}

	type MockResultCase struct {
		matches       []models.User
		errGetMatches error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			matches: []models.User{
				{
					ID:          1,
					Name:        "Drip))",
					Email:       "drip@app.ru",
					Password:    "hah",
					Date:        "2001-02-22",
					Age:         19,
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
					Age:         20,
					Description: "em privet",
					Imgs:        []string{"4", "3"},
					Tags:        []string{"JS"},
				},
			},
			errGetMatches: nil,
		},
		// Test ErrContextNilError
		{
			matches:       []models.User{},
			errGetMatches: nil,
		},
		// Test ErrorNotFound
		{
			matches:       []models.User{},
			errGetMatches: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.user.ID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, testCase.user))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetUsersMatchesWithSearching",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("string")).Return(MockResultCases[i].matches, MockResultCases[i].errGetMatches)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2, nil)
		allMatches, status := testUserUsecase.UsersMatchesWithSearching(r.Context(), testCase.searchingData)

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(testCase.matches, allMatches)

	}
}

func TestUserUsecase_UserLikes(t *testing.T) {
	type TestCase struct {
		user  models.User
		likes models.Likes
		err   error
	}
	testCases := []TestCase{
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
			likes: models.Likes{
				AllUsers: map[uint64]models.User{
					0: {
						ID:          1,
						Name:        "Владимир",
						Email:       "drip@app.ru",
						Password:    "hah",
						Date:        "2001-02-22",
						Age:         19,
						Description: "vsem privet",
						Imgs:        []string{"4", "3"},
						Tags:        []string{"BMSTU"},
					},
					1: {
						ID:          2,
						Name:        "влаД",
						Email:       "dr@app.ru",
						Password:    "hah",
						Date:        "2000-02-22",
						Age:         20,
						Description: "em privet",
						Imgs:        []string{"4", "3"},
						Tags:        []string{"JS"},
					},
				},
				Count: "2",
			},
			err: nil,
		},
		// Test ErrContextNilError
		{
			user: models.User{
				ID: 2,
			},
			err: models.ErrContextNilError,
		},
		// Test ErrorNotFound
		{
			user:  models.User{},
			likes: models.Likes{},
			err:   errors.New(""),
		},
	}

	type MockResultCase struct {
		likes         []models.User
		errGetMatches error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			likes: []models.User{
				{
					ID:          1,
					Name:        "Drip))",
					Email:       "drip@app.ru",
					Password:    "hah",
					Date:        "2001-02-22",
					Age:         19,
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
					Age:         20,
					Description: "em privet",
					Imgs:        []string{"4", "3"},
					Tags:        []string{"JS"},
				},
			},
			errGetMatches: nil,
		},
		// Test ErrContextNilError
		{
			likes:         []models.User{},
			errGetMatches: nil,
		},
		// Test ErrorNotFound
		{
			likes:         []models.User{},
			errGetMatches: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.user.ID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, testCase.user))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetUsersLikes",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].likes, MockResultCases[i].errGetMatches)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2, nil)
		allLikes, status := testUserUsecase.UserLikes(r.Context())

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(testCase.likes, allLikes)

	}
}

func TestUserUsecase_GetAllReports(t *testing.T) {
	type TestCase struct {
		user    models.User
		reports models.Reports
		err     error
	}
	testCases := []TestCase{
		// Test OK
		{
			user: models.User{},
			reports: models.Reports{
				AllReports: map[uint64]models.Report{
					0: {ReportDesc: "anime"},
					1: {ReportDesc: "BMSTU"},
					2: {ReportDesc: "walk"},
					3: {ReportDesc: "netflix"},
					4: {ReportDesc: "prikolchiki"},
				},
				Count: 5,
			},
			err: nil,
		},
		// Test ErrGetTags
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
			reports: models.Reports{},
			err:     errors.New(""),
		},
	}

	type MockResultCase struct {
		reports       map[uint64]string
		errGetReports error
	}
	MockResultCases := []MockResultCase{
		// Test OK
		{
			reports: map[uint64]string{
				0: "anime",
				1: "BMSTU",
				2: "walk",
				3: "netflix",
				4: "prikolchiki",
			},
			errGetReports: nil,
		},
		// Test ErrGetTags
		{
			reports:       map[uint64]string{},
			errGetReports: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.user.ID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, testCase.user))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("GetReports",
			mock.AnythingOfType("*context.timerCtx")).Return(MockResultCases[i].reports, MockResultCases[i].errGetReports)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2, nil)
		allReports, status := testUserUsecase.GetAllReports(r.Context())

		assert.Equal(t, testCase.err, status, message)
		reflect.DeepEqual(testCase.reports, allReports)

	}
}

func TestUserUsecase_AddReport(t *testing.T) {
	type TestCase struct {
		user   models.User
		report models.NewReport
		err    error
	}
	testCases := []TestCase{
		// Test OK FAKE
		{
			user: models.User{},
			report: models.NewReport{
				ToId:       0,
				ReportDesc: models.FakeReport,
			},
			err: nil,
		},
		// Test OK AGGRESSION
		{
			user: models.User{},
			report: models.NewReport{
				ToId:       0,
				ReportDesc: models.AggressionReport,
			},
			err: nil,
		},
		// Test OK SKAM
		{
			user: models.User{},
			report: models.NewReport{
				ToId:       0,
				ReportDesc: models.SkamReport,
			},
			err: nil,
		},
		// Test OK UNDERAGE
		{
			user: models.User{},
			report: models.NewReport{
				ToId:       0,
				ReportDesc: models.UnderageReport,
			},
			err: nil,
		},
		// Test ErrContextNilError
		{
			user: models.User{
				ID: 2,
			},
			report: models.NewReport{},
			err:    models.ErrContextNilError,
		},
		// Test ErrAddReport
		{
			user:   models.User{},
			report: models.NewReport{},
			err:    errors.New(""),
		},
		// Test ErrDeleteReaction
		{
			user:   models.User{},
			report: models.NewReport{},
			err:    errors.New(""),
		},
		// Test ErrDeleteMatches
		{
			user:   models.User{},
			report: models.NewReport{},
			err:    errors.New(""),
		},
		// Test ErrAddReaction
		{
			user:   models.User{},
			report: models.NewReport{},
			err:    errors.New(""),
		},
		// Test ErrGetReportsCount
		{
			user:   models.User{},
			report: models.NewReport{},
			err:    errors.New(""),
		},
		// Test ErrGetReportsWithMaxCount
		{
			user:   models.User{},
			report: models.NewReport{},
			err:    errors.New(""),
		},
		// Test ErrGetReportDesc
		{
			user:   models.User{},
			report: models.NewReport{},
			err:    errors.New(""),
		},
		// Test ErrUpdateReportStatus
		{
			user:   models.User{},
			report: models.NewReport{},
			err:    errors.New(""),
		},
	}

	type MockResultCase struct {
		reportsCount              uint64
		banId                     uint64
		banDesc                   string
		errAddReport              error
		errDeleteReaction         error
		errDeleteMatches          error
		errAddReaction            error
		errGetReportsCount        error
		errGetReportsWithMaxCount error
		errGetReportDesc          error
		errUpdateReportStatus     error
	}
	MockResultCases := []MockResultCase{
		// Test OK FAKE
		{
			reportsCount: 5,
			banId:        2,
			banDesc:      models.FakeReport,
		},
		// Test OK AGGRESSION
		{
			reportsCount: 5,
			banId:        2,
			banDesc:      models.AggressionReport,
		},
		// Test OK SKAM
		{
			reportsCount: 5,
			banId:        2,
			banDesc:      models.SkamReport,
		},
		// Test OK UNDERAGE
		{
			reportsCount: 5,
			banId:        2,
			banDesc:      models.UnderageReport,
		},
		// Test ErrContextNilError
		{
			reportsCount: 5,
			banId:        2,
			banDesc:      "",
		},
		// Test ErrAddReport
		{
			reportsCount: 5,
			banId:        2,
			banDesc:      "",
			errAddReport: errors.New(""),
		},
		// Test ErrDeleteReaction
		{
			reportsCount:      5,
			banId:             2,
			banDesc:           "",
			errDeleteReaction: errors.New(""),
		},
		// Test ErrDeleteMatches
		{
			reportsCount:     5,
			banId:            2,
			banDesc:          "",
			errDeleteMatches: errors.New(""),
		},
		// Test ErrAddReaction
		{
			reportsCount:   5,
			banId:          2,
			banDesc:        "",
			errAddReaction: errors.New(""),
		},
		// Test ErrGetReportsCount
		{
			reportsCount:       5,
			banId:              2,
			banDesc:            "",
			errGetReportsCount: errors.New(""),
		},
		// Test ErrGetReportsWithMaxCount
		{
			reportsCount:              5,
			banId:                     2,
			banDesc:                   "",
			errGetReportsWithMaxCount: errors.New(""),
		},
		// Test ErrGetReportDesc
		{
			reportsCount:     5,
			banId:            2,
			banDesc:          "",
			errGetReportDesc: errors.New(""),
		},
		// Test ErrUpdateReportStatus
		{
			reportsCount:          5,
			banId:                 2,
			banDesc:               "",
			errUpdateReportStatus: errors.New(""),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		r, err := http.NewRequest(http.MethodGet, "test", nil)
		assert.NoError(t, err)
		if testCase.user.ID != 2 {
			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, testCase.user))
		}

		mockUserRepository := new(userMocks.UserRepository)
		mockUserRepository.On("AddReport",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("models.NewReport")).Return(MockResultCases[i].errAddReport)
		mockUserRepository.On("DeleteReaction",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].errDeleteReaction)
		mockUserRepository.On("DeleteMatches",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].errDeleteMatches)
		mockUserRepository.On("AddReaction",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].errAddReaction)
		mockUserRepository.On("GetReportsCount",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].reportsCount, MockResultCases[i].errGetReportsCount)
		mockUserRepository.On("GetReportsWithMaxCount",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].banId, MockResultCases[i].errGetReportsWithMaxCount)
		mockUserRepository.On("GetReportDesc",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64")).Return(MockResultCases[i].banDesc, MockResultCases[i].errGetReportDesc)
		mockUserRepository.On("UpdateReportStatus",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("uint64"),
			mock.AnythingOfType("string")).Return(MockResultCases[i].errUpdateReportStatus)
		mockFileRepository := new(fileMocks.FileRepository)

		testUserUsecase := usecase.NewUserUsecase(mockUserRepository, mockFileRepository, time.Second*2, nil)
		err = testUserUsecase.AddReport(r.Context(), testCase.report)

		assert.Equal(t, testCase.err, err, message)

	}
}
