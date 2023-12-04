package app

import (
	"github.com/jessicatarra/greenlight/ms/auth/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreate(t *testing.T) {
	// Initialize the repository mock
	userRepo := mocks.UserRepository{}
	tokenRepo := mocks.TokenRepository{}

	// Create the app instance with the repository mock
	app := NewApp(&userRepo, &tokenRepo)

	// Prepare the input for the Create function
	input := CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	// Set the expectations on the user repository mock
	userRepo.On("InsertNewUser", mock.AnythingOfType("*entity.User")).Return(nil)
	tokenRepo.On("New", mock.Anything, mock.AnythingOfType("time.Duration"), mock.IsType("string")).Return(nil, nil)

	// Call the Create function
	user, err := app.Create(input)

	// Assert the results
	assert.NotNil(t, user)
	assert.NoError(t, err)
}
