// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// WS is an autogenerated mock type for the WS type
type WS struct {
	mock.Mock
}

// ReadJSON provides a mock function with given fields: _a0
func (_m *WS) ReadJSON(_a0 interface{}) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteJSON provides a mock function with given fields: _a0
func (_m *WS) WriteJSON(_a0 interface{}) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}