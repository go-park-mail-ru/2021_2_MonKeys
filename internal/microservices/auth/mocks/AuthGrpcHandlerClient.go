// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"
	auth "dripapp/internal/microservices/auth/delivery/grpc/protobuff"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// AuthGrpcHandlerClient is an autogenerated mock type for the AuthGrpcHandlerClient type
type AuthGrpcHandlerClient struct {
	mock.Mock
}

// GetById provides a mock function with given fields: ctx, in, opts
func (_m *AuthGrpcHandlerClient) GetById(ctx context.Context, in *auth.Session, opts ...grpc.CallOption) (*auth.User, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *auth.User
	if rf, ok := ret.Get(0).(func(context.Context, *auth.Session, ...grpc.CallOption) *auth.User); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *auth.Session, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFromSession provides a mock function with given fields: ctx, in, opts
func (_m *AuthGrpcHandlerClient) GetFromSession(ctx context.Context, in *auth.Cookie, opts ...grpc.CallOption) (*auth.Session, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *auth.Session
	if rf, ok := ret.Get(0).(func(context.Context, *auth.Cookie, ...grpc.CallOption) *auth.Session); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *auth.Cookie, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
