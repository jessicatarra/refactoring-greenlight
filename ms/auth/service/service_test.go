//go:build integration
// +build integration

package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jessicatarra/greenlight/ms/auth/app"
	"github.com/jessicatarra/greenlight/ms/auth/app/mocks"
	"github.com/jessicatarra/greenlight/ms/auth/entity"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResource_CreateUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		//Arrange
		router := httprouter.New()

		mockApp := &mocks.App{}

		res := resource{app: mockApp}

		RegisterHandlers(mockApp, router)

		requestBody := []byte(`{
				"name": "John Doe",
				"email": "johndoe@example.com",
				"password": "password123"
			}`)

		req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		resRec := httptest.NewRecorder()

		expectedInput := app.CreateUserRequest{
			Name:     "John Doe",
			Email:    "johndoe@example.com",
			Password: "password123",
		}

		expectedUser := &entity.User{
			ID:    1,
			Name:  "John Doe",
			Email: "johndoe@example.com",
		}

		mockApp.On("Create", expectedInput).Return(expectedUser, nil)

		// Act
		res.create(resRec, req)

		//Assert
		if resRec.Code != http.StatusCreated {
			t.Errorf("unexpected status code: got %d, want %d", resRec.Code, http.StatusCreated)
		}

		var responseBody map[string]*entity.User
		err := json.Unmarshal(resRec.Body.Bytes(), &responseBody)
		if err != nil {
			t.Errorf("failed to parse response body: %s", err)
		}

		if responseBody["user"] == nil {
			t.Errorf("expected 'user' field in response body, got nil")
		} else {
			// Check the user entity fields
			if responseBody["user"].ID != expectedUser.ID {
				t.Errorf("unexpected user ID: got %d, want %d", responseBody["user"].ID, expectedUser.ID)
			}

			if responseBody["user"].Name != expectedUser.Name {
				t.Errorf("unexpected user name: got %s, want %s", responseBody["user"].Name, expectedUser.Name)
			}

			if responseBody["user"].Email != expectedUser.Email {
				t.Errorf("unexpected user email: got %s, want %s", responseBody["user"].Email, expectedUser.Email)
			}
		}

		mockApp.AssertCalled(t, "Create", expectedInput)
	})

	t.Run("error", func(t *testing.T) {
		//Arrange
		router := httprouter.New()

		mockApp := &mocks.App{}

		res := resource{app: mockApp}

		RegisterHandlers(mockApp, router)

		requestBody := []byte(`{
				"name": "John Doe",
				"email": "johndoe@example.com",
				"password": "password123"
			}`)

		req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		resRec := httptest.NewRecorder()

		expectedInput := app.CreateUserRequest{
			Name:     "John Doe",
			Email:    "johndoe@example.com",
			Password: "password123",
		}

		expectedErr := errors.New("The server encountered a problem and could not process your request")

		mockApp.On("Create", expectedInput).Return(nil, expectedErr)

		//Act
		res.create(resRec, req)

		if resRec.Code != http.StatusInternalServerError {
			t.Errorf("unexpected status code: got %d, want %d", resRec.Code, http.StatusInternalServerError)
		}

		var responseBody map[string]interface{}
		err := json.Unmarshal(resRec.Body.Bytes(), &responseBody)
		if err != nil {
			t.Errorf("failed to parse response body: %s", err)
		}

		if responseBody["Error"] == nil {
			t.Errorf("expected 'error' field in response body, got nil")
		} else {
			// Check the error message
			if responseBody["Error"] != expectedErr.Error() {
				t.Errorf("unexpected error message: got %s, want %s", responseBody["Error"], expectedErr.Error())
			}
		}

		mockApp.AssertCalled(t, "Create", expectedInput)
	})
}
