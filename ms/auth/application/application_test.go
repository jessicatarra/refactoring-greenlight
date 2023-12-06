//go:build integration
// +build integration

package application

import (
	"errors"
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/internal/jsonlog"
	"github.com/jessicatarra/greenlight/ms/auth/domain"
	"github.com/jessicatarra/greenlight/ms/auth/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
)

func Init() (mocks.UserRepository, mocks.TokenRepository, *jsonlog.Logger, *sync.WaitGroup, config.Config) {
	userRepo := mocks.UserRepository{}
	tokenRepo := mocks.TokenRepository{}
	logger := &jsonlog.Logger{}
	wg := &sync.WaitGroup{}
	cfg := config.Config{
		Smtp: struct {
			Host     string
			Port     int
			Username string
			Password string
			Sender   string
		}{
			Host:     "sandbox.smtp.mailtrap.io",
			Port:     25,
			Username: "username",
			Password: "password",
			Sender:   "Greenlight <no-reply@tarralva.com>",
		},
	}
	return userRepo, tokenRepo, logger, wg, cfg
}

func TestApp_CreateUseCase(t *testing.T) {
	// Initialize the repositories mock
	userRepo, tokenRepo, logger, wg, cfg := Init()

	// CreateUseCase the application instance with the repositories mock
	app := NewAppl(&userRepo, &tokenRepo, logger, wg, cfg)

	// Prepare the input for the CreateUseCase function
	input := domain.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	userRepo.On("InsertNewUser", mock.AnythingOfType("*domain.User")).Return(nil)
	tokenRepo.On("New", mock.Anything, mock.AnythingOfType("time.Duration"), mock.IsType("string")).Return(nil, nil)

	// Call the CreateUseCase function
	user, err := app.CreateUseCase(input)

	// Assert the results
	assert.NotNil(t, user)
	assert.NoError(t, err)
}

func TestApp_GetByEmailUseCase(t *testing.T) {
	// Initialize the repositories mock
	userRepo, tokenRepo, logger, wg, cfg := Init()

	// CreateUseCase the application instance with the repositories mock
	app := NewAppl(&userRepo, &tokenRepo, logger, wg, cfg)

	// Prepare the input for the CreateUseCase function
	input := domain.CreateUserRequest{
		Name:     "Sarah Foo",
		Email:    "sarah@example.com",
		Password: "password123",
	}

	userRepo.On("GetUserByEmail", mock.AnythingOfType("string")).Return(nil, errors.New("record not found"))
	tokenRepo.On("New", mock.Anything, mock.AnythingOfType("time.Duration"), mock.IsType("string")).Return(nil, nil)

	// Call the CreateUseCase function
	user, err := app.GetByEmailUseCase(input)

	// Assert the results
	assert.Error(t, err)
	assert.Nil(t, user)
}
