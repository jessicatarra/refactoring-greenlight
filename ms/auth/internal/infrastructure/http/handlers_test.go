//go:build auth
// +build auth

package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jessicatarra/greenlight/internal/password"
	"github.com/jessicatarra/greenlight/ms/auth/internal/domain"
	"github.com/jessicatarra/greenlight/ms/auth/internal/domain/mocks"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouterAndMocks() (*mocks.Appl, Handlers) {
	mockApp := &mocks.Appl{}

	res := registerHandlers(mockApp)

	return mockApp, res
}

func createRequestBody() []byte {
	requestBody := []byte(`{
		"name": "John Doe",
		"email": "johndoe@example.com",
		"password": "password123"
	}`)

	return requestBody
}

func assertStatusCode(t *testing.T, resRec *httptest.ResponseRecorder, expectedStatusCode int) {
	if resRec.Code != expectedStatusCode {
		t.Errorf("unexpected status code: got %d, want %d", resRec.Code, expectedStatusCode)
	}
}

func assertResponseBody(t *testing.T, resRec *httptest.ResponseRecorder, responseBody interface{}) {
	err := json.Unmarshal(resRec.Body.Bytes(), responseBody)
	if err != nil {
		t.Errorf("failed to parse response body: %s", err)
	}
}

func assertUserFields(t *testing.T, responseBody map[string]*domain.User, expectedUser *domain.User) {
	if responseBody["user"] == nil {
		t.Errorf("expected 'user' field in response body, got nil")
	} else {
		user := responseBody["user"]
		if user.ID != expectedUser.ID {
			t.Errorf("unexpected user ID: got %d, want %d", user.ID, expectedUser.ID)
		}
		if user.Name != expectedUser.Name {
			t.Errorf("unexpected user name: got %s, want %s", user.Name, expectedUser.Name)
		}
		if user.Email != expectedUser.Email {
			t.Errorf("unexpected user email: got %s, want %s", user.Email, expectedUser.Email)
		}
	}
}

func TestResource_Create(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockApp, res := setupRouterAndMocks()
		expectedInput := &domain.CreateUserRequest{
			Name:     "John Doe",
			Email:    "johndoe@example.com",
			Password: "password123",
		}
		expectedUser := &domain.User{
			ID:    1,
			Name:  "John Doe",
			Email: "johndoe@example.com",
		}

		requestBody := createRequestBody()
		req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resRec := httptest.NewRecorder()

		// Mock CreateUseCase and GetByEmailUseCase
		mockApp.On("CreateUseCase", expectedInput, mock.AnythingOfType("string")).Return(expectedUser, nil)
		mockApp.On("GetByEmailUseCase", "johndoe@example.com").Return(nil, errors.New("record not found"))

		// Act
		res.createUser(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusCreated)
		var responseBody map[string]*domain.User
		assertResponseBody(t, resRec, &responseBody)
		assertUserFields(t, responseBody, expectedUser)
	})

	t.Run("error - GetByEmailUseCase return error", func(t *testing.T) {
		// Arrange
		mockApp, res := setupRouterAndMocks()
		expectedInput := &domain.CreateUserRequest{
			Name:     "John Doe",
			Email:    "johndoe@example.com",
			Password: "password123",
		}
		expectedUser := &domain.User{
			ID:    1,
			Name:  "John Doe",
			Email: "johndoe@example.com",
		}

		requestBody := createRequestBody()
		req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resRec := httptest.NewRecorder()

		// Mock CreateUseCase and GetByEmailUseCase
		mockApp.On("GetByEmailUseCase", expectedUser.Email).Return(nil, errors.New("error"))
		mockApp.On("CreateUseCase", expectedInput, mock.AnythingOfType("string")).Return(expectedUser, nil)

		// Act
		res.createUser(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusInternalServerError)
	})

	t.Run("error - bad request", func(t *testing.T) {
		// Arrange
		mockApp, res := setupRouterAndMocks()
		expectedInput := &domain.CreateUserRequest{
			Name:     "John Doe",
			Email:    "johndoe@example.com",
			Password: "password123",
		}
		expectedUser := &domain.User{
			ID:    1,
			Name:  "John Doe",
			Email: "johndoe@example.com",
		}

		requestBody := []byte(`{
		"email": "johndoe@example.com",
		"password": "password123"
	}`)
		req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resRec := httptest.NewRecorder()

		// Mock CreateUseCase and GetByEmailUseCase
		mockApp.On("CreateUseCase", expectedInput, mock.AnythingOfType("string")).Return(expectedUser, nil)
		mockApp.On("GetByEmailUseCase", "johndoe@example.com").Return(nil, errors.New("record not found"))

		// Act
		res.createUser(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusUnprocessableEntity)
	})

	t.Run("error - validate user", func(t *testing.T) {
		// Arrange
		mockApp, res := setupRouterAndMocks()
		expectedInput := &domain.CreateUserRequest{
			Name:     "John Doe",
			Email:    "johndoe@example.com",
			Password: "password123",
		}
		expectedUser := &domain.User{
			ID:    1,
			Name:  "John Doe",
			Email: "johndoe@example.com",
		}

		req := httptest.NewRequest(http.MethodPost, "/v1/users", nil)
		req.Header.Set("Content-Type", "application/json")
		resRec := httptest.NewRecorder()

		// Mock CreateUseCase and GetByEmailUseCase
		mockApp.On("CreateUseCase", expectedInput, mock.AnythingOfType("string")).Return(expectedUser, nil)
		mockApp.On("GetByEmailUseCase", "johndoe@example.com").Return(nil, errors.New("record not found"))

		// Act
		res.createUser(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusBadRequest)
	})

	t.Run("error", func(t *testing.T) {
		// Arrange
		mockApp, res := setupRouterAndMocks()
		expectedInput := &domain.CreateUserRequest{
			Name:     "John Doe",
			Email:    "johndoe@example.com",
			Password: "password123",
		}

		requestBody := createRequestBody()
		req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resRec := httptest.NewRecorder()
		expectedErr := errors.New("The server encountered a problem and could not process your request")

		// Mock CreateUseCase and GetByEmailUseCase
		mockApp.On("CreateUseCase", expectedInput, mock.AnythingOfType("string")).Return(nil, expectedErr)
		mockApp.On("GetByEmailUseCase", "johndoe@example.com").Return(nil, errors.New("record not found"))

		// Act
		res.createUser(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusInternalServerError)
		var responseBody map[string]string
		assertResponseBody(t, resRec, &responseBody)
		if responseBody["Error"] == "" {
			t.Errorf("expected 'error' fieldin response body, got nil")
		} else {
			if responseBody["Error"] != expectedErr.Error() {
				t.Errorf("unexpected error message: got %s, want %s", responseBody["Error"], expectedErr.Error())
			}
		}
	})
}

func TestResource_Activate(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockApp, res := setupRouterAndMocks()
		activationToken := "GQRPVONORIEUPDJ6V4RTDIVSTQ"
		expectedInput := domain.ActivateUserRequest{
			TokenPlaintext: activationToken,
		}
		expectedUser := &domain.User{
			ID:        1,
			Name:      "John Doe",
			Email:     "johndoe@example.com",
			Activated: true,
		}

		req := httptest.NewRequest(http.MethodPut, "/v1/users/activate/?token="+activationToken, nil)
		resRec := httptest.NewRecorder()

		// Mock ActivateUseCase
		mockApp.On("ActivateUseCase", expectedInput.TokenPlaintext).Return(expectedUser, nil)

		// Act
		res.activateUser(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusCreated)
		var responseBody map[string]*domain.User
		assertResponseBody(t, resRec, &responseBody)
		assertUserFields(t, responseBody, expectedUser)
		mockApp.AssertCalled(t, "ActivateUseCase", expectedInput.TokenPlaintext)
	})

	t.Run("error", func(t *testing.T) {
		// Arrange
		mockApp, res := setupRouterAndMocks()
		activationToken := "GQRPVONORIEUPDJ6V4RTDIVSTQ"
		expectedInput := domain.ActivateUserRequest{
			TokenPlaintext: activationToken,
		}

		req := httptest.NewRequest(http.MethodPut, "/v1/users/activate/?token="+activationToken, nil)
		resRec := httptest.NewRecorder()
		expectedErr := errors.New("The server encountered a problem and could not process your request")

		// Mock ActivateUseCase
		mockApp.On("ActivateUseCase", expectedInput.TokenPlaintext).Return(nil, expectedErr)

		// Act
		res.activateUser(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusInternalServerError)
		var responseBody map[string]interface{}
		assertResponseBody(t, resRec, &responseBody)
		if responseBody["Error"] == nil {
			t.Errorf("expected 'error' field in response body, got nil")
		} else {
			if responseBody["Error"] != expectedErr.Error() {
				t.Errorf("unexpected error message: got %s, want %s", responseBody["Error"], expectedErr.Error())
			}
		}
		mockApp.AssertCalled(t, "ActivateUseCase", expectedInput.TokenPlaintext)
	})
}

func TestResource_AuthenticationToken(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockApp, res := setupRouterAndMocks()
		hashedPassword, _ := password.Hash("password123")
		expectedUser := &domain.User{
			ID:             1,
			Name:           "John Doe",
			Email:          "johndoe@example.com",
			HashedPassword: hashedPassword,
			Activated:      true,
		}

		requestBody := []byte(`{
		"email": "johndoe@example.com",
		"password": "password123"
		}`)

		req := httptest.NewRequest(http.MethodPost, "/v1/tokens/authentication", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resRec := httptest.NewRecorder()

		// Mock GetByEmailUseCase and CreateAuthTokenUseCase
		mockApp.On("GetByEmailUseCase", expectedUser.Email).Return(expectedUser, nil)
		mockApp.On("CreateAuthTokenUseCase", expectedUser.ID).Return([]byte("thisisasecreT"), nil)

		// Act
		res.createAuthenticationToken(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusCreated)
		var responseBody map[string]string
		assertResponseBody(t, resRec, &responseBody)
	})

	t.Run("error - bad request status code", func(t *testing.T) {
		// Arrange
		mockApp, res := setupRouterAndMocks()
		hashedPassword, _ := password.Hash("password123")
		expectedUser := &domain.User{
			ID:             1,
			Name:           "John Doe",
			Email:          "johndoe@example.com",
			HashedPassword: hashedPassword,
			Activated:      true,
		}
		requestBody := []byte(`{
		"email": "johndoe@example.com",
		}`)

		req := httptest.NewRequest(http.MethodPost, "/v1/tokens/authentication", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resRec := httptest.NewRecorder()

		// Mock GetByEmailUseCase and CreateAuthTokenUseCase
		mockApp.On("GetByEmailUseCase", expectedUser.Email).Return(expectedUser, nil)
		mockApp.On("CreateAuthTokenUseCase", expectedUser.ID).Return([]byte("thisisasecreT"), nil)

		// Act
		res.createAuthenticationToken(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusBadRequest)
	})

	t.Run("error - GetByEmailUseCase return domain.ErrRecordNotFound error", func(t *testing.T) {
		// Arrange
		mockApp, res := setupRouterAndMocks()
		hashedPassword, _ := password.Hash("password123")
		expectedUser := &domain.User{
			ID:             1,
			Name:           "John Doe",
			Email:          "johndoe@example.com",
			HashedPassword: hashedPassword,
			Activated:      true,
		}
		requestBody := []byte(`{
		"email": "johndoe@example.com",
		"password": "password123"
		}`)

		req := httptest.NewRequest(http.MethodPost, "/v1/tokens/authentication", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resRec := httptest.NewRecorder()

		// Mock GetByEmailUseCase and CreateAuthTokenUseCase
		mockApp.On("GetByEmailUseCase", expectedUser.Email).Return(nil, domain.ErrRecordNotFound)
		mockApp.On("CreateAuthTokenUseCase", expectedUser.ID).Return([]byte("thisisasecreT"), nil)

		// Act
		res.createAuthenticationToken(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusUnauthorized)
	})

	t.Run("error - GetByEmailUseCase return error", func(t *testing.T) {
		// Arrange
		mockApp, res := setupRouterAndMocks()
		hashedPassword, _ := password.Hash("password123")
		expectedUser := &domain.User{
			ID:             1,
			Name:           "John Doe",
			Email:          "johndoe@example.com",
			HashedPassword: hashedPassword,
			Activated:      true,
		}
		requestBody := []byte(`{
		"email": "johndoe@example.com",
		"password": "password123"
		}`)

		req := httptest.NewRequest(http.MethodPost, "/v1/tokens/authentication", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resRec := httptest.NewRecorder()

		// Mock GetByEmailUseCase and CreateAuthTokenUseCase
		mockApp.On("GetByEmailUseCase", expectedUser.Email).Return(nil, errors.New("error"))
		mockApp.On("CreateAuthTokenUseCase", expectedUser.ID).Return([]byte("thisisasecreT"), nil)

		// Act
		res.createAuthenticationToken(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusInternalServerError)
	})

	t.Run("error - password matches return error", func(t *testing.T) {
		// Arrange
		mockApp, res := setupRouterAndMocks()
		expectedUser := &domain.User{
			ID:        1,
			Name:      "John Doe",
			Email:     "johndoe@example.com",
			Activated: true,
		}
		requestBody := []byte(`{
		"email": "johndoe@example.com",
		"password": "password123"
		}`)

		req := httptest.NewRequest(http.MethodPost, "/v1/tokens/authentication", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resRec := httptest.NewRecorder()

		// Mock GetByEmailUseCase and CreateAuthTokenUseCase
		mockApp.On("GetByEmailUseCase", expectedUser.Email).Return(expectedUser, nil)
		mockApp.On("CreateAuthTokenUseCase", expectedUser.ID).Return([]byte("thisisasecreT"), nil)

		// Act
		res.createAuthenticationToken(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusInternalServerError)
	})

	t.Run("error - input validator has error", func(t *testing.T) {
		// Arrange
		hashedPassword, _ := password.Hash("password123")

		mockApp, res := setupRouterAndMocks()
		expectedUser := &domain.User{
			ID:             1,
			Name:           "John Doe",
			Email:          "johndoe@example.com",
			HashedPassword: hashedPassword,
			Activated:      true,
		}
		requestBody := []byte(`{
		"email": "johndoe@example.com",
		"password": "password1"
		}`)

		req := httptest.NewRequest(http.MethodPost, "/v1/tokens/authentication", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resRec := httptest.NewRecorder()

		// Mock GetByEmailUseCase and CreateAuthTokenUseCase
		mockApp.On("GetByEmailUseCase", expectedUser.Email).Return(expectedUser, nil)
		mockApp.On("CreateAuthTokenUseCase", expectedUser.ID).Return([]byte("thisisasecreT"), nil)

		// Act
		res.createAuthenticationToken(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusUnprocessableEntity)
	})

	t.Run("error - token return error", func(t *testing.T) {
		// Arrange
		mockApp, res := setupRouterAndMocks()
		hashedPassword, _ := password.Hash("password123")
		expectedUser := &domain.User{
			ID:             1,
			Name:           "John Doe",
			Email:          "johndoe@example.com",
			HashedPassword: hashedPassword,
			Activated:      true,
		}

		requestBody := []byte(`{
		"email": "johndoe@example.com",
		"password": "password123"
		}`)

		req := httptest.NewRequest(http.MethodPost, "/v1/tokens/authentication", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resRec := httptest.NewRecorder()

		// Mock GetByEmailUseCase and CreateAuthTokenUseCase
		mockApp.On("GetByEmailUseCase", expectedUser.Email).Return(expectedUser, nil)
		mockApp.On("CreateAuthTokenUseCase", expectedUser.ID).Return(nil, errors.New("error"))

		// Act
		res.createAuthenticationToken(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusInternalServerError)
	})
}
