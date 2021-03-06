// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	models "dripapp/internal/dripapp/models"
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// FileRepository is an autogenerated mock type for the FileRepository type
type FileRepository struct {
	mock.Mock
}

// CreateFoldersForNewUser provides a mock function with given fields: user
func (_m *FileRepository) CreateFoldersForNewUser(user models.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(models.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: filePath
func (_m *FileRepository) Delete(filePath string) error {
	ret := _m.Called(filePath)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(filePath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveUserPhoto provides a mock function with given fields: user, file, fileName
func (_m *FileRepository) SaveUserPhoto(user models.User, file io.Reader, fileName string) (string, error) {
	ret := _m.Called(user, file, fileName)

	var r0 string
	if rf, ok := ret.Get(0).(func(models.User, io.Reader, string) string); ok {
		r0 = rf(user, file, fileName)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.User, io.Reader, string) error); ok {
		r1 = rf(user, file, fileName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
