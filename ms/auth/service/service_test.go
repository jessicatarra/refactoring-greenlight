//go:build integration
// +build integration

package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jessicatarra/greenlight/internal/utils/helpers"
	"github.com/jessicatarra/greenlight/ms/auth/domain"
	"github.com/jessicatarra/greenlight/ms/auth/domain/mocks"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouterAndMocks() (*mocks.Appl, resource) {
	router := httprouter.New()
	mockApp := &mocks.Appl{}
	mockHelpers := helpers.New()
	res := resource{appl: mockApp, helpers: mockHelpers}
	RegisterHandlers(mockApp, router)
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
		expectedInput := domain.CreateUserRequest{
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
		mockApp.On("CreateUseCase", expectedInput).Return(expectedUser, nil)
		mockApp.On("GetByEmailUseCase", "johndoe@example.com").Return(nil, errors.New("record not found"))

		// Act
		res.createUser(resRec, req)

		// Assert
		assertStatusCode(t, resRec, http.StatusCreated)
		var responseBody map[string]*domain.User
		assertResponseBody(t, resRec, &responseBody)
		assertUserFields(t, responseBody, expectedUser)
		mockApp.AssertCalled(t, "CreateUseCase", expectedInput)
	})

	t.Run("error", func(t *testing.T) {
		// Arrange
		mockApp, res := setupRouterAndMocks()
		expectedInput := domain.CreateUserRequest{
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
		mockApp.On("CreateUseCase", expectedInput).Return(nil, expectedErr)
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

		mockApp.AssertCalled(t, "CreateUseCase", expectedInput)

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
