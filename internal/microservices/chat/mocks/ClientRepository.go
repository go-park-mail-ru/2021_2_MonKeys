// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	models "dripapp/internal/microservices/chat/models"

	mock "github.com/stretchr/testify/mock"
)

// ClientRepository is an autogenerated mock type for the ClientRepository type
type ClientRepository struct {
	mock.Mock
}

// SaveMessage provides a mock function with given fields: userId, toId, text
func (_m *ClientRepository) SaveMessage(userId uint64, toId uint64, text string) (models.Message, error) {
	ret := _m.Called(userId, toId, text)

	var r0 models.Message
	if rf, ok := ret.Get(0).(func(uint64, uint64, string) models.Message); ok {
		r0 = rf(userId, toId, text)
	} else {
		r0 = ret.Get(0).(models.Message)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, uint64, string) error); ok {
		r1 = rf(userId, toId, text)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}