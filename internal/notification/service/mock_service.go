// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package service

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// AddPushToken provides a mock function with given fields: ctx, req
func (_m *MockService) AddPushToken(ctx context.Context, req AddPushTokenReq) error {
	ret := _m.Called(ctx, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, AddPushTokenReq) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemovePushToken provides a mock function with given fields: ctx, req
func (_m *MockService) RemovePushToken(ctx context.Context, req RemovePushTokenReq) error {
	ret := _m.Called(ctx, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, RemovePushTokenReq) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendNotification provides a mock function with given fields: ctx, req
func (_m *MockService) SendNotification(ctx context.Context, req SendNotificationReq) error {
	ret := _m.Called(ctx, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, SendNotificationReq) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}