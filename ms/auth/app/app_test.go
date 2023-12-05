package app

import (
	"github.com/jessicatarra/greenlight/internal/validator"
	"github.com/jessicatarra/greenlight/ms/auth/entity"
	"github.com/jessicatarra/greenlight/ms/auth/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestValidateUser(t *testing.T) {
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

func TestValidateEmail(t *testing.T) {
	v := validator.New()
	email := "john@example.com"

	ValidateEmail(v, email)

	assert.True(t, v.Valid())
}

func TestValidatePasswordPlaintext(t *testing.T) {
	v := validator.New()
	password := "password123"

	ValidatePasswordPlaintext(v, password)

	assert.True(t, v.Valid())
}

func TestCreate(t *testing.T) {
	// Initialize the repositories mock
	userRepo := mocks.UserRepository{}
	tokenRepo := mocks.TokenRepository{}

	// Create the app instance with the repositories mock
	app := NewApp(&userRepo, &tokenRepo)

	// Prepare the input for the Create function
	input := CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	// Set the expectations on the user repositories mock
	userRepo.On("InsertNewUser", mock.AnythingOfType("*entity.User")).Return(nil)
	tokenRepo.On("New", mock.Anything, mock.AnythingOfType("time.Duration"), mock.IsType("string")).Return(nil, nil)

	// Call the Create function
	user, err := app.Create(input)

	// Assert the results
	assert.NotNil(t, user)
	assert.NoError(t, err)
}
