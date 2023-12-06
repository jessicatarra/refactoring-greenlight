// Code generated by mockery v2.33.2. DO NOT EDIT.

package mocks

import (
	domain "github.com/jessicatarra/greenlight/ms/auth/domain"
	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// GetForToken provides a mock function with given fields: tokenScope, tokenPlaintext
func (_m *UserRepository) GetForToken(tokenScope string, tokenPlaintext string) (*domain.User, error) {
	ret := _m.Called(tokenScope, tokenPlaintext)

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*domain.User, error)); ok {
		return rf(tokenScope, tokenPlaintext)
	}
	if rf, ok := ret.Get(0).(func(string, string) *domain.User); ok {
		r0 = rf(tokenScope, tokenPlaintext)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(tokenScope, tokenPlaintext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByEmail provides a mock function with given fields: email
func (_m *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
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

// GetUserById provides a mock function with given fields: id
func (_m *UserRepository) GetUserById(id int64) (*domain.User, error) {
	ret := _m.Called(id)

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (*domain.User, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) *domain.User); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertNewUser provides a mock function with given fields: user
func (_m *UserRepository) InsertNewUser(user *domain.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateUser provides a mock function with given fields: user
func (_m *UserRepository) UpdateUser(user *domain.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
