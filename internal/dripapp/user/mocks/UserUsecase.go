// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"
	io "io"

	mock "github.com/stretchr/testify/mock"

	models "dripapp/internal/dripapp/models"
)

// UserUsecase is an autogenerated mock type for the UserUsecase type
type UserUsecase struct {
	mock.Mock
}

// AddPhoto provides a mock function with given fields: c, photo
func (_m *UserUsecase) AddPhoto(c context.Context, photo io.Reader) (string, models.HTTPError) {
	ret := _m.Called(c, photo)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, io.Reader) string); ok {
		r0 = rf(c, photo)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 models.HTTPError
	if rf, ok := ret.Get(1).(func(context.Context, io.Reader) models.HTTPError); ok {
		r1 = rf(c, photo)
	} else {
		r1 = ret.Get(1).(models.HTTPError)
	}

	return r0, r1
}

// CurrentUser provides a mock function with given fields: c
func (_m *UserUsecase) CurrentUser(c context.Context) (models.User, models.HTTPError) {
	ret := _m.Called(c)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(context.Context) models.User); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 models.HTTPError
	if rf, ok := ret.Get(1).(func(context.Context) models.HTTPError); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Get(1).(models.HTTPError)
	}

	return r0, r1
}

// DeletePhoto provides a mock function with given fields: c, photo
func (_m *UserUsecase) DeletePhoto(c context.Context, photo models.Photo) models.HTTPError {
	ret := _m.Called(c, photo)

	var r0 models.HTTPError
	if rf, ok := ret.Get(0).(func(context.Context, models.Photo) models.HTTPError); ok {
		r0 = rf(c, photo)
	} else {
		r0 = ret.Get(0).(models.HTTPError)
	}

	return r0
}

// EditProfile provides a mock function with given fields: c, newUserData
func (_m *UserUsecase) EditProfile(c context.Context, newUserData models.User) (models.User, models.HTTPError) {
	ret := _m.Called(c, newUserData)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(context.Context, models.User) models.User); ok {
		r0 = rf(c, newUserData)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 models.HTTPError
	if rf, ok := ret.Get(1).(func(context.Context, models.User) models.HTTPError); ok {
		r1 = rf(c, newUserData)
	} else {
		r1 = ret.Get(1).(models.HTTPError)
	}

	return r0, r1
}

// GetAllTags provides a mock function with given fields: c
func (_m *UserUsecase) GetAllTags(c context.Context) (models.Tags, models.HTTPError) {
	ret := _m.Called(c)

	var r0 models.Tags
	if rf, ok := ret.Get(0).(func(context.Context) models.Tags); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Get(0).(models.Tags)
	}

	var r1 models.HTTPError
	if rf, ok := ret.Get(1).(func(context.Context) models.HTTPError); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Get(1).(models.HTTPError)
	}

	return r0, r1
}

// Login provides a mock function with given fields: c, logUserData
func (_m *UserUsecase) Login(c context.Context, logUserData models.LoginUser) (models.User, models.HTTPError) {
	ret := _m.Called(c, logUserData)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(context.Context, models.LoginUser) models.User); ok {
		r0 = rf(c, logUserData)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 models.HTTPError
	if rf, ok := ret.Get(1).(func(context.Context, models.LoginUser) models.HTTPError); ok {
		r1 = rf(c, logUserData)
	} else {
		r1 = ret.Get(1).(models.HTTPError)
	}

	return r0, r1
}

// NextUser provides a mock function with given fields: c
func (_m *UserUsecase) NextUser(c context.Context) ([]models.User, models.HTTPError) {
	ret := _m.Called(c)

	var r0 []models.User
	if rf, ok := ret.Get(0).(func(context.Context) []models.User); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.User)
		}
	}

	var r1 models.HTTPError
	if rf, ok := ret.Get(1).(func(context.Context) models.HTTPError); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Get(1).(models.HTTPError)
	}

	return r0, r1
}

// Reaction provides a mock function with given fields: c, reactionData
func (_m *UserUsecase) Reaction(c context.Context, reactionData models.UserReaction) (models.Match, models.HTTPError) {
	ret := _m.Called(c, reactionData)

	var r0 models.Match
	if rf, ok := ret.Get(0).(func(context.Context, models.UserReaction) models.Match); ok {
		r0 = rf(c, reactionData)
	} else {
		r0 = ret.Get(0).(models.Match)
	}

	var r1 models.HTTPError
	if rf, ok := ret.Get(1).(func(context.Context, models.UserReaction) models.HTTPError); ok {
		r1 = rf(c, reactionData)
	} else {
		r1 = ret.Get(1).(models.HTTPError)
	}

	return r0, r1
}

// Signup provides a mock function with given fields: c, logUserData
func (_m *UserUsecase) Signup(c context.Context, logUserData models.LoginUser) (models.User, models.HTTPError) {
	ret := _m.Called(c, logUserData)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(context.Context, models.LoginUser) models.User); ok {
		r0 = rf(c, logUserData)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 models.HTTPError
	if rf, ok := ret.Get(1).(func(context.Context, models.LoginUser) models.HTTPError); ok {
		r1 = rf(c, logUserData)
	} else {
		r1 = ret.Get(1).(models.HTTPError)
	}

	return r0, r1
}

// UsersMatches provides a mock function with given fields: c
func (_m *UserUsecase) UsersMatches(c context.Context) (models.Matches, models.HTTPError) {
	ret := _m.Called(c)

	var r0 models.Matches
	if rf, ok := ret.Get(0).(func(context.Context) models.Matches); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Get(0).(models.Matches)
	}

	var r1 models.HTTPError
	if rf, ok := ret.Get(1).(func(context.Context) models.HTTPError); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Get(1).(models.HTTPError)
	}

	return r0, r1
}
