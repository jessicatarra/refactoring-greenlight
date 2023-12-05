//go:build integration
// +build integration

package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jessicatarra/greenlight/ms/auth/entity"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInsertNewUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)

	hash := []byte("sampleHash")

	t.Run("Success", func(t *testing.T) {
		// Arrange
		user := &entity.User{
			Name:      "John Doe",
			Email:     "johndoe@example.com",
			Password:  entity.Password{Hash: hash},
			Activated: true,
		}

		rows := sqlmock.NewRows([]string{"id", "created_at", "version"}).
			AddRow(1, time.Now(), 1)

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Name, user.Email, user.Password.Hash, user.Activated).WillReturnRows(rows)

		// Act
		err := repo.InsertNewUser(user)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		user := &entity.User{
			Name:      "John Doe",
			Email:     "johndoe@example.com",
			Password:  entity.Password{Hash: hash},
			Activated: true,
		}

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Name, user.Email, user.Password.Hash, user.Activated).
			WillReturnError(errors.New("some error"))

		// Act
		err := repo.InsertNewUser(user)

		// Assert
		assert.Error(t, err)
	})
}

func TestGetUserByEmail(t *testing.T) {
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

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	hash := []byte("sampleHash")

	t.Run("Success", func(t *testing.T) {
		// Arrange
		user := &entity.User{
			ID:        1,
			Name:      "John Doe",
			Email:     "johndoe@example.com",
			Password:  entity.Password{Hash: hash},
			Activated: true,
			Version:   1,
		}

		rows := sqlmock.NewRows([]string{"version"}).
			AddRow(2)

		mock.ExpectQuery("UPDATE users").WithArgs(user.Name, user.Email, user.Password.Hash, user.Activated, user.ID, user.Version).WillReturnRows(rows)

		// Act
		err := repo.UpdateUser(user)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		user := &entity.User{
			ID:        1,
			Name:      "John Doe",
			Email:     "johndoe@example.com",
			Password:  entity.Password{Hash: hash},
			Activated: true,
			Version:   1,
		}

		mock.ExpectExec("UPDATE users").
			WithArgs(user.Name, user.Email, user.Password.Hash, user.Activated, user.Version, user.ID).
			WillReturnError(errors.New("some error"))

		// Act
		err := repo.UpdateUser(user)

		// Assert
		assert.Error(t, err)
	})
}

func TestGetUserById(t *testing.T) {
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
