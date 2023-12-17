//go:build auth
// +build auth

package application

import (
	"errors"
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/internal/password"
	"github.com/jessicatarra/greenlight/ms/auth/internal/domain"
	"github.com/jessicatarra/greenlight/ms/auth/internal/domain/mocks"
	"github.com/jessicatarra/greenlight/ms/auth/internal/infrastructure/repositories"
	"github.com/pascaldekloe/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strconv"
	"sync"
	"testing"
)

func Init() (mocks.UserRepository, mocks.TokenRepository, mocks.PermissionRepository, config.Config, sync.WaitGroup) {
	userRepo := mocks.UserRepository{}
	tokenRepo := mocks.TokenRepository{}
	permissionRepo := mocks.PermissionRepository{}
	wg := sync.WaitGroup{}
	cfg := config.Config{
		Jwt: struct {
			Secret string
		}{
			Secret: "ifTp39TukiePBVu7SY1K+l07v8l1aiP+F2Tu9BxQ34c=",
		},
		Smtp: struct {
			Host     string
			Port     int
			Username string
			Password string
			From     string
		}{
			Host:     "sandbox.smtp.mailtrap.io",
			Port:     25,
			Username: "username",
			Password: "password",
			From:     "Greenlight <no-reply@tarralva.com>",
		},
		Auth: struct {
			HttpBaseURL    string
			GrpcBaseURL    string
			GrpcServerPort int
			HttpPort       int
		}{
			HttpBaseURL:    "localhost:8082",
			GrpcBaseURL:    "localhost:50051",
			GrpcServerPort: 50051,
			HttpPort:       8082,
		},
	}
	return userRepo, tokenRepo, permissionRepo, cfg, wg
}

func TestAppl_CreateUseCase(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		// Initialize the repositories mock
		userRepo, tokenRepo, permissionRepo, cfg, wg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)

		// Prepare the input for the CreateUseCase function
		input := domain.CreateUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		hashedPassword, _ := password.Hash(input.Password)

		// Set up the success step
		userRepo.On("InsertNewUser", mock.AnythingOfType("*domain.User"), mock.AnythingOfType("string")).Return(nil)
		permissionRepo.On("AddForUser", mock.AnythingOfTypeArgument("int64"), "movies:read").Return(nil)
		tokenRepo.On("New", mock.Anything, mock.AnythingOfType("time.Duration"), mock.IsType("string")).Return(nil, nil)

		// Call the CreateUseCase function
		user, err := app.CreateUseCase(&input, hashedPassword)

		// Assert the results
		assert.NotNil(t, user)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		// Initialize the repositories mock
		userRepo, tokenRepo, permissionRepo, cfg, wg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)

		// Prepare the input for the CreateUseCase function
		input := domain.CreateUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}
		hashedPassword, _ := password.Hash(input.Password)

		// Set up the error step
		userRepo.On("InsertNewUser", mock.AnythingOfType("*domain.User"), mock.AnythingOfType("string")).Return(errors.New("failed to insert user"))
		tokenRepo.On("New", mock.Anything, mock.AnythingOfType("time.Duration"), mock.IsType("string")).Return(nil, nil)

		// Call the CreateUseCase function again
		user, err := app.CreateUseCase(&input, hashedPassword)

		// Assert the error step
		assert.Nil(t, user)
		assert.Error(t, err)
	})

}

func TestAppl_GetByEmailUseCase(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Initialize the repositories mock
		userRepo, tokenRepo, permissionRepo, cfg, wg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)

		// Prepare the input for the CreateUseCase function
		input := domain.CreateUserRequest{
			Name:     "Sarah Foo",
			Email:    "sarah@example.com",
			Password: "password123",
		}

		userRepo.On("GetUserByEmail", mock.AnythingOfType("string")).Return(nil, errors.New("record not found"))
		tokenRepo.On("New", mock.Anything, mock.AnythingOfType("time.Duration"), mock.IsType("string")).Return(nil, nil)

		// Call the CreateUseCase function
		user, err := app.GetByEmailUseCase(input.Email)

		// Assert the results
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("error", func(t *testing.T) {
		// Initialize the repositories mock
		userRepo, tokenRepo, permissionRepo, cfg, wg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)

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
		user, err := app.GetByEmailUseCase(input.Email)

		// Assert the error step
		assert.Nil(t, user)
		assert.Error(t, err)
	})
}

func TestAppl_ActivateUseCase(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		// Initialize the repositories mock
		userRepo, tokenRepo, permissionRepo, cfg, wg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)

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
		userRepo, tokenRepo, permissionRepo, cfg, wg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)

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
		userRepo, tokenRepo, permissionRepo, cfg, wg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)

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
		userRepo, tokenRepo, permissionRepo, cfg, wg := Init()

		// CreateUseCase the application instance with the repositories mock
		app := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)

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

func TestAppl_CreateAuthTokenUseCase(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		// Arrange
		userRepo, tokenRepo, permissionRepo, cfg, wg := Init()
		appl := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)

		expectedUserID := int64(1)
		expectedSubject := strconv.FormatInt(expectedUserID, 10)
		expectedIssuer := cfg.Auth.HttpBaseURL
		expectedAudience := []string{cfg.Auth.HttpBaseURL}

		// Act
		tokenBytes, err := appl.CreateAuthTokenUseCase(expectedUserID)
		token, err := jwt.HMACCheck(tokenBytes, []byte(cfg.Jwt.Secret))

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, token.Subject, expectedSubject)
		assert.Equal(t, token.Issuer, expectedIssuer)
		assert.Equal(t, token.Audiences, expectedAudience)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		cfg := config.Config{
			Auth: struct {
				HttpBaseURL    string
				GrpcBaseURL    string
				GrpcServerPort int
				HttpPort       int
			}{
				HttpBaseURL:    "localhost:8082",
				GrpcBaseURL:    "localhost:50051",
				GrpcServerPort: 50051,
				HttpPort:       8082,
			},
		}
		userRepo, tokenRepo, permissionRepo, _, wg := Init()
		appl := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)
		expectedUserID := int64(1)

		// Act
		_, err := appl.CreateAuthTokenUseCase(expectedUserID)

		// Assert
		assert.Error(t, err)
	})
}

func TestAppl_ValidateAuthTokenUseCase(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		userRepo, tokenRepo, permissionRepo, cfg, wg := Init()
		appl := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)
		expectedUserID := int64(1)
		expectedUser := &domain.User{
			ID:        int64(1),
			Name:      "John Doe",
			Email:     "john@example.com",
			Activated: true,
		}
		userRepo.On("GetUserById", mock.AnythingOfType("int64")).Return(expectedUser, nil)

		// Act
		tokenBytes, err := appl.CreateAuthTokenUseCase(expectedUserID)
		token := string(tokenBytes)
		user, err := appl.ValidateAuthTokenUseCase(token)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, user.ID, expectedUserID)
	})

	t.Run("Error - JWT Secret", func(t *testing.T) {
		// Arrange
		userRepo, tokenRepo, permissionRepo, _, wg := Init()
		cfg := config.Config{
			Auth: struct {
				HttpBaseURL    string
				GrpcBaseURL    string
				GrpcServerPort int
				HttpPort       int
			}{
				HttpBaseURL:    "localhost:8082",
				GrpcBaseURL:    "localhost:50051",
				GrpcServerPort: 50051,
				HttpPort:       8082,
			},
		}
		appl := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)
		expectedUserID := int64(1)
		expectedUser := &domain.User{
			ID:        int64(1),
			Name:      "John Doe",
			Email:     "john@example.com",
			Activated: true,
		}
		userRepo.On("GetUserById", mock.AnythingOfType("int64")).Return(expectedUser, nil)

		// Act
		tokenBytes, err := appl.CreateAuthTokenUseCase(expectedUserID)
		token := string(tokenBytes)
		_, err = appl.ValidateAuthTokenUseCase(token)

		// Assert
		assert.Error(t, err)
	})

	t.Run("Error - database", func(t *testing.T) {
		// Arrange
		userRepo, tokenRepo, permissionRepo, cfg, wg := Init()
		appl := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)
		expectedUserID := int64(1)
		userRepo.On("GetUserById", mock.AnythingOfType("int64")).Return(nil, errors.New("record not found"))

		// Act
		tokenBytes, err := appl.CreateAuthTokenUseCase(expectedUserID)
		token := string(tokenBytes)
		_, err = appl.ValidateAuthTokenUseCase(token)

		// Assert
		assert.Error(t, err)
	})

}

func TestAppl_UserPermissionUseCase(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		userRepo, tokenRepo, permissionRepo, cfg, wg := Init()
		appl := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)

		expectedUserID := int64(1)
		code := "movie:read"
		permissions := domain.Permissions{code}

		permissionRepo.On("GetAllForUser", mock.AnythingOfType("int64")).Return(permissions, nil)

		// Act
		err := appl.UserPermissionUseCase(code, expectedUserID)

		// Assert
		assert.NoError(t, err)
	})
	t.Run("Error - database", func(t *testing.T) {
		// Arrange
		userRepo, tokenRepo, permissionRepo, cfg, wg := Init()
		appl := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)

		expectedUserID := int64(1)
		code := "movie:read"

		permissionRepo.On("GetAllForUser", mock.AnythingOfType("int64")).Return(nil, errors.New("error"))

		// Act
		err := appl.UserPermissionUseCase(code, expectedUserID)

		// Assert
		assert.Error(t, err)
	})
	t.Run("Error - permission not included", func(t *testing.T) {
		// Arrange
		userRepo, tokenRepo, permissionRepo, cfg, wg := Init()
		appl := NewAppl(&userRepo, &tokenRepo, &permissionRepo, &wg, cfg)

		expectedUserID := int64(1)
		code := "movie:read"
		permissions := domain.Permissions{code}

		permissionRepo.On("GetAllForUser", mock.AnythingOfType("int64")).Return(permissions, nil)

		// Act
		err := appl.UserPermissionUseCase("movies:write", expectedUserID)

		// Assert
		assert.Error(t, err)
	})
}
