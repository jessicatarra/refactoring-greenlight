// Code generated by mockery v2.33.2. DO NOT EDIT.

package mocks

import (
	http "net/http"

	httprouter "github.com/julienschmidt/httprouter"
	mock "github.com/stretchr/testify/mock"

	service "github.com/jessicatarra/greenlight/ms/auth/internal/service"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Handlers provides a mock function with given fields: router
func (_m *Service) Handlers(router *httprouter.Router) {
	_m.Called(router)
}

// Middlewares provides a mock function with given fields:
func (_m *Service) Middlewares() service.Middlewares {
	ret := _m.Called()

	var r0 service.Middlewares
	if rf, ok := ret.Get(0).(func() service.Middlewares); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.Middlewares)
		}
	}

	return r0
}

// Routes provides a mock function with given fields:
func (_m *Service) Routes() http.Handler {
	ret := _m.Called()

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func() http.Handler); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
