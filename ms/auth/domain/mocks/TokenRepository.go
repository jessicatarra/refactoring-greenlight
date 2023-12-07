// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/jessicatarra/greenlight/ms/auth/domain"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// TokenRepository is an autogenerated mock type for the TokenRepository type
type TokenRepository struct {
	mock.Mock
}

// DeleteAllForUser provides a mock function with given fields: scope, userID
func (_m *TokenRepository) DeleteAllForUser(scope string, userID int64) error {
	ret := _m.Called(scope, userID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAllForUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int64) error); ok {
		r0 = rf(scope, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Insert provides a mock function with given fields: token
func (_m *TokenRepository) Insert(token *domain.Token) error {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for Insert")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Token) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// New provides a mock function with given fields: userID, ttl, scope
func (_m *TokenRepository) New(userID int64, ttl time.Duration, scope string) (*domain.Token, error) {
	ret := _m.Called(userID, ttl, scope)

	if len(ret) == 0 {
		panic("no return value specified for New")
	}

	var r0 *domain.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, time.Duration, string) (*domain.Token, error)); ok {
		return rf(userID, ttl, scope)
	}
	if rf, ok := ret.Get(0).(func(int64, time.Duration, string) *domain.Token); ok {
		r0 = rf(userID, ttl, scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, time.Duration, string) error); ok {
		r1 = rf(userID, ttl, scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTokenRepository creates a new instance of TokenRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTokenRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *TokenRepository {
	mock := &TokenRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
