//go:build integration
// +build integration

package application

import (
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/internal/jsonlog"
	"github.com/jessicatarra/greenlight/internal/validator"
	"github.com/jessicatarra/greenlight/ms/auth/entity"
	"github.com/jessicatarra/greenlight/ms/auth/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
)

func TestApp_ValidateUser(t *testing.T) {
	v := validator.New()

	password := "password123!"
	hash := []byte("sampleHash")

	user := &entity.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: entity.Password{Plaintext: &password, Hash: hash},
	}

	ValidateUser(v, user)

	assert.True(t, v.Valid())
}

func TestApp_ValidateEmail(t *testing.T) {
	v := validator.New()
	email := "john@example.com"

	ValidateEmail(v, email)

	assert.True(t, v.Valid())
}

func TestApp_ValidatePasswordPlaintext(t *testing.T) {
	v := validator.New()
	password := "password123"

	ValidatePasswordPlaintext(v, password)

	assert.True(t, v.Valid())
}

func TestApp_Create(t *testing.T) {
	// Initialize the repositories mock
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

	// CreateUseCase the application instance with the repositories mock
	app := NewAppl(&userRepo, &tokenRepo, logger, wg, cfg)

	// Prepare the input for the CreateUseCase function
	input := entity.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	userRepo.On("InsertNewUser", mock.AnythingOfType("*entity.User")).Return(nil)
	tokenRepo.On("New", mock.Anything, mock.AnythingOfType("time.Duration"), mock.IsType("string")).Return(nil, nil)

	// Call the CreateUseCase function
	user, err := app.CreateUseCase(input)

	// Assert the results
	assert.NotNil(t, user)
	assert.NoError(t, err)
}
