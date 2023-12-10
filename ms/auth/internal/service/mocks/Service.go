// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	domain "github.com/jessicatarra/greenlight/ms/auth/internal/domain"

	httprouter "github.com/julienschmidt/httprouter"

	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Routes provides a mock function with given fields:
func (_m *Service) Routes() http.Handler {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Routes")
	}

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

// logRequestMiddleware provides a mock function with given fields: next
func (_m *Service) logRequestMiddleware(next http.Handler) http.Handler {
	ret := _m.Called(next)

	if len(ret) == 0 {
		panic("no return value specified for logRequestMiddleware")
	}

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func(http.Handler) http.Handler); ok {
		r0 = rf(next)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// registerHandlers provides a mock function with given fields: appl, router
func (_m *Service) registerHandlers(appl domain.Appl, router *httprouter.Router) {
	_m.Called(appl, router)
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
