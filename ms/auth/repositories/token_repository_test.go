//go:build auth
// +build auth

package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jessicatarra/greenlight/ms/auth/domain"
	"github.com/jessicatarra/greenlight/ms/auth/domain/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTokenRepository_New(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockTokenInterface := mocks.TokenInterface{}

	repo := &tokenRepository{
		db:    db,
		token: &mockTokenInterface,
	}

	// Test cases
	t.Run("Success", func(t *testing.T) {
		// Arrange
		userID := int64(1)
		ttl := 1 * time.Hour
		scope := ScopeActivation

		expectedToken := &domain.Token{
			Plaintext: "mock_token",
			Hash:      []byte("mock_hash"),
			UserID:    userID,
			Expiry:    time.Now().Add(ttl),
			Scope:     scope,
		}
		mock.ExpectExec("INSERT INTO tokens").
			WithArgs(sqlmock.AnyArg(), userID, sqlmock.AnyArg(), scope).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockTokenInterface.On("GenerateToken", userID, ttl, scope).Return(expectedToken, nil)

		// Act
		newToken, err := repo.New(userID, ttl, scope)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, newToken)
		assert.Equal(t, expectedToken, newToken)
	})
}

func TestTokenRepository_NewError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockTokenInterface := mocks.TokenInterface{}

	repo := &tokenRepository{
		db:    db,
		token: &mockTokenInterface,
	}

	t.Run("Error", func(t *testing.T) {
		// Arrange
		userID := int64(1)
		ttl := 1 * time.Hour
		scope := ScopeActivation

		mock.ExpectExec("INSERT INTO tokens").
			WithArgs(sqlmock.AnyArg(), userID, ttl, scope).
			WillReturnError(sqlmock.ErrCancelled)
		mockTokenInterface.On("GenerateToken", userID, ttl, scope).Return(nil, errors.New("some error"))

		// Act
		newToken, err := repo.New(userID, ttl, scope)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, newToken)
	})
}

func TestTokenRepository_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTokenRepo(db)
	hash := []byte("sampleHash")

	// Test cases
	t.Run("Success", func(t *testing.T) {
		// Arrange
		expectedToken := &domain.Token{
			Hash:   hash,
			UserID: 1,
			Expiry: time.Now().Add(1 * time.Hour),
			Scope:  ScopeActivation,
		}

		mock.ExpectExec("INSERT INTO tokens").
			WithArgs(expectedToken.Hash, expectedToken.UserID, expectedToken.Expiry, expectedToken.Scope).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		err := repo.Insert(expectedToken)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		token := &domain.Token{
			Hash:   hash,
			UserID: 1,
			Expiry: time.Now().Add(1 * time.Hour),
			Scope:  ScopeActivation,
		}

		mock.ExpectExec("INSERT INTO tokens").
			WithArgs(token.Hash, token.UserID, token.Expiry, token.Scope).
			WillReturnError(sqlmock.ErrCancelled)

		// Act
		err := repo.Insert(token)

		// Assert
		assert.Error(t, err)
	})
}

func TestTokenRepository_DeleteAllForUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTokenRepo(db)

	// Test cases
	t.Run("Success", func(t *testing.T) {
		// Arrange
		scope := ScopeActivation
		userID := int64(1)

		mock.ExpectExec("DELETE FROM tokens").
			WithArgs(scope, userID).
			WillReturnResult(sqlmock.NewResult(0, 2))

		// Act
		err := repo.DeleteAllForUser(scope, userID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		scope := ScopeActivation
		userID := int64(1)

		mock.ExpectExec("DELETE FROM tokens").
			WithArgs(scope, userID).
			WillReturnError(sqlmock.ErrCancelled)

		// Act
		err := repo.DeleteAllForUser(scope, userID)

		// Assert
		assert.Error(t, err)
	})
}

func TestTokenRepository_GetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)

	// Test cases
	t.Run("Success", func(t *testing.T) {
		// Arrange
		userID := int64(1)

		rows := sqlmock.NewRows([]string{"id", "created_at", "name", "email", "password_hash", "activated", "version"}).
			AddRow(userID, time.Now(), "John Doe", "johndoe@example.com", "somehash", true, 1)

		mock.ExpectQuery("SELECT").
			WithArgs(userID).
			WillReturnRows(rows)

		// Act
		user, err := repo.GetUserById(userID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		userID := int64(1)

		mock.ExpectQuery("SELECT").
			WithArgs(userID).
			WillReturnError(sqlmock.ErrCancelled)

		// Act
		user, err := repo.GetUserById(userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
