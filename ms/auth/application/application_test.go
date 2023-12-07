//go:build integration
// +build integration

package application

import (
	"errors"
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/ms/auth/domain"
	"github.com/jessicatarra/greenlight/ms/auth/domain/mocks"
	"github.com/jessicatarra/greenlight/ms/auth/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Init() (mocks.UserRepository, mocks.TokenRepository, config.Config) {
	userRepo := mocks.UserRepository{}
	tokenRepo := mocks.TokenRepository{}
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
	return userRepo, tokenRepo, cfg
}

func TestApp_CreateUseCase(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		// Initialize the repositories mock
		userRepo, tokenRepo, cfg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, cfg)

		// Prepare the input for the CreateUseCase function
		input := domain.CreateUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		// Set up the success step
		userRepo.On("InsertNewUser", mock.AnythingOfType("*domain.User")).Return(nil)
		tokenRepo.On("New", mock.Anything, mock.AnythingOfType("time.Duration"), mock.IsType("string")).Return(nil, nil)

		// Call the CreateUseCase function
		user, err := app.CreateUseCase(input)

		// Assert the results
		assert.NotNil(t, user)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		// Initialize the repositories mock
		userRepo, tokenRepo, cfg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, cfg)

		// Prepare the input for the CreateUseCase function
		input := domain.CreateUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		// Set up the error step
		userRepo.On("InsertNewUser", mock.AnythingOfType("*domain.User")).Return(errors.New("failed to insert user"))
		tokenRepo.On("New", mock.Anything, mock.AnythingOfType("time.Duration"), mock.IsType("string")).Return(nil, nil)

		// Call the CreateUseCase function again
		user, err := app.CreateUseCase(input)

		// Assert the error step
		assert.Nil(t, user)
		assert.Error(t, err)
	})

}

func TestApp_GetByEmailUseCase(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Initialize the repositories mock
		userRepo, tokenRepo, cfg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, cfg)

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
	})

	t.Run("error", func(t *testing.T) {
		// Initialize the repositories mock
		userRepo, tokenRepo, cfg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, cfg)

		// Prepare the input for the CreateUseCase function
		input := domain.CreateUserRequest{
			Name:     "Sarah Foo",
			Email:    "sarah@example.com",
			Password: "password123",
		}

		// Set up the error step
		userRepo.On("GetUserByEmail", mock.AnythingOfType("string")).Return(nil, errors.New("database error"))
		tokenRepo.On("New", mock.Anything, mock.AnythingOfType("time.Duration"), mock.IsType("string")).Return(nil, nil)

		// Call the GetByEmailUseCase function again
		user, err := app.GetByEmailUseCase(input)

		// Assert the error step
		assert.Nil(t, user)
		assert.Error(t, err)
	})
}

func TestApp_ActivateUseCase(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		// Initialize the repositories mock
		userRepo, tokenRepo, cfg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, cfg)

		// Prepare the input for the ActivateUseCase function
		tokenPlainText := "valid_token"

		// Set up the success step
		expectedUser := &domain.User{
			ID:        int64(1),
			Name:      "John Doe",
			Email:     "john@example.com",
			Activated: true,
		}
		userRepo.On("GetForToken", repositories.ScopeActivation, tokenPlainText).Return(expectedUser, nil)
		userRepo.On("UpdateUser", expectedUser).Return(nil)
		tokenRepo.On("DeleteAllForUser", repositories.ScopeActivation, expectedUser.ID).Return(nil)

		// Call the ActivateUseCase function
		user, err := app.ActivateUseCase(tokenPlainText)

		// Assert the success step
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, err)
	})

	t.Run("error - GetForToken", func(t *testing.T) {
		// Initialize the repositories mock
		userRepo, tokenRepo, cfg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, cfg)

		// Prepare the input for the ActivateUseCase function
		tokenPlainText := "valid_token"

		// Set up the error step for GetForToken
		expectedErr := errors.New("failed to get user for token")
		userRepo.On("GetForToken", repositories.ScopeActivation, tokenPlainText).Return(nil, expectedErr)

		// Call the ActivateUseCase function again
		user, err := app.ActivateUseCase(tokenPlainText)

		// Assert the error step for GetForToken
		assert.Nil(t, user)
		assert.Equal(t, err.Error(), expectedErr.Error())
	})

	t.Run("error - UpdateUser", func(t *testing.T) {
		// Initialize the repositories mock
		userRepo, tokenRepo, cfg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, cfg)

		// Prepare the input for the ActivateUseCase function
		tokenPlainText := "valid_token"

		// Set up the error step for UpdateUser
		expectedUser := &domain.User{
			ID:        int64(1),
			Name:      "John Doe",
			Email:     "john@example.com",
			Activated: true,
		}
		expectedErr := errors.New("failed to update user")
		userRepo.On("GetForToken", repositories.ScopeActivation, tokenPlainText).Return(expectedUser, nil)
		userRepo.On("UpdateUser", expectedUser).Return(expectedErr)

		// Call the ActivateUseCase function again
		user, err := app.ActivateUseCase(tokenPlainText)

		// Assert the error step for UpdateUser
		assert.Nil(t, user)
		assert.Equal(t, err.Error(), expectedErr.Error())
	})

	t.Run("error - DeleteAllForUser", func(t *testing.T) {
		// Initialize the repositories mock
		userRepo, tokenRepo, cfg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, cfg)

		// Prepare the input for the ActivateUseCase function
		tokenPlainText := "valid_token"

		// Set up the error step for DeleteAllForUser
		expectedUser := &domain.User{
			ID:        int64(1),
			Name:      "John Doe",
			Email:     "john@example.com",
			Activated: true,
		}
		expectedErr := errors.New("failed to delete tokens")
		userRepo.On("GetForToken", repositories.ScopeActivation, tokenPlainText).Return(expectedUser, nil)
		userRepo.On("UpdateUser", expectedUser).Return(nil)
		tokenRepo.On("DeleteAllForUser", repositories.ScopeActivation, expectedUser.ID).Return(expectedErr)

		// Call the ActivateUseCase function again
		user, err := app.ActivateUseCase(tokenPlainText)

		// Assert the error step for DeleteAllForUser
		assert.Nil(t, user)
		assert.Equal(t, err.Error(), expectedErr.Error())
	})
}
