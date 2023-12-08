//go:build auth
// +build auth

package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jessicatarra/greenlight/internal/password"
	"github.com/jessicatarra/greenlight/ms/auth/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserRepository_InsertNewUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)

	hashedPassword, _ := password.Hash("password123!")

	t.Run("Success", func(t *testing.T) {
		// Arrange
		user := &domain.User{
			Name:           "John Doe",
			Email:          "johndoe@example.com",
			HashedPassword: hashedPassword,
			Activated:      true,
		}

		rows := sqlmock.NewRows([]string{"id", "created_at", "version"}).
			AddRow(1, time.Now(), 1)

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Name, user.Email, user.HashedPassword, user.Activated).WillReturnRows(rows)

		// Act
		err := repo.InsertNewUser(user, hashedPassword)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		user := &domain.User{
			Name:           "John Doe",
			Email:          "johndoe@example.com",
			HashedPassword: hashedPassword,
			Activated:      true,
		}

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Name, user.Email, user.HashedPassword, user.Activated).
			WillReturnError(errors.New("some error"))

		// Act
		err := repo.InsertNewUser(user, hashedPassword)

		// Assert
		assert.Error(t, err)
	})
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	hash := []byte("sampleHash")

	t.Run("Success", func(t *testing.T) {
		// Arrange
		email := "johndoe@example.com"

		rows := sqlmock.NewRows([]string{"id", "created_at", "name", "email", "password_hash", "activated", "version"}).
			AddRow(1, time.Now(), "John Doe", "johndoe@example.com", hash, true, 1)

		mock.ExpectQuery("SELECT").
			WithArgs(email).
			WillReturnRows(rows)

		// Act
		user, err := repo.GetUserByEmail(email)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		email := "johndoe@example.com"

		mock.ExpectQuery("SELECT").
			WithArgs(email).
			WillReturnError(errors.New("some error"))

		// Act
		user, err := repo.GetUserByEmail(email)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_UpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	hashedPassword, _ := password.Hash("password123!")

	t.Run("Success", func(t *testing.T) {
		// Arrange
		user := &domain.User{
			ID:             1,
			Name:           "John Doe",
			Email:          "johndoe@example.com",
			HashedPassword: hashedPassword,
			Activated:      true,
			Version:        1,
		}

		rows := sqlmock.NewRows([]string{"version"}).
			AddRow(2)

		mock.ExpectQuery("UPDATE users").WithArgs(user.Name, user.Email, user.HashedPassword, user.Activated, user.ID, user.Version).WillReturnRows(rows)

		// Act
		err := repo.UpdateUser(user)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		user := &domain.User{
			ID:             1,
			Name:           "John Doe",
			Email:          "johndoe@example.com",
			HashedPassword: hashedPassword,
			Activated:      true,
			Version:        1,
		}

		mock.ExpectExec("UPDATE users").
			WithArgs(user.Name, user.Email, user.HashedPassword, user.Activated, user.Version, user.ID).
			WillReturnError(errors.New("some error"))

		// Act
		err := repo.UpdateUser(user)

		// Assert
		assert.Error(t, err)
	})
}

func TestUserRepository_GetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)

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
			WillReturnError(errors.New("some error"))

		// Act
		user, err := repo.GetUserById(userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
