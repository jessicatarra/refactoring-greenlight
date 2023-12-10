// Code generated by mockery v2.33.2. DO NOT EDIT.

package mocks

import (
	domain "github.com/jessicatarra/greenlight/ms/auth/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// Appl is an autogenerated mock type for the Appl type
type Appl struct {
	mock.Mock
}

// ActivateUseCase provides a mock function with given fields: tokenPlainText
func (_m *Appl) ActivateUseCase(tokenPlainText string) (*domain.User, error) {
	ret := _m.Called(tokenPlainText)

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.User, error)); ok {
		return rf(tokenPlainText)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.User); ok {
		r0 = rf(tokenPlainText)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tokenPlainText)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAuthTokenUseCase provides a mock function with given fields: userID
func (_m *Appl) CreateAuthTokenUseCase(userID int64) ([]byte, error) {
	ret := _m.Called(userID)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) ([]byte, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(int64) []byte); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUseCase provides a mock function with given fields: input, hashedPassword
func (_m *Appl) CreateUseCase(input domain.CreateUserRequest, hashedPassword string) (*domain.User, error) {
	ret := _m.Called(input, hashedPassword)

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.CreateUserRequest, string) (*domain.User, error)); ok {
		return rf(input, hashedPassword)
	}
	if rf, ok := ret.Get(0).(func(domain.CreateUserRequest, string) *domain.User); ok {
		r0 = rf(input, hashedPassword)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.CreateUserRequest, string) error); ok {
		r1 = rf(input, hashedPassword)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByEmailUseCase provides a mock function with given fields: email
func (_m *Appl) GetByEmailUseCase(email string) (*domain.User, error) {
	ret := _m.Called(email)

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.User, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.User); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAppl creates a new instance of Appl. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAppl(t interface {
	mock.TestingT
	Cleanup(func())
}) *Appl {
	mock := &Appl{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}