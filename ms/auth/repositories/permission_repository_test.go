package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPermissionRepository_GetAllForUser(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPermissionRepo(db)

	t.Run("Success", func(t *testing.T) {
		// Arrange
		userID := int64(1)
		rows := sqlmock.NewRows([]string{"codes"}).AddRow("movies:read")
		mock.ExpectQuery("SELECT").
			WithArgs(userID).
			WillReturnRows(rows)

		// Act
		permissions, err := repo.GetAllForUser(userID)
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, permissions)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		userID := int64(1)
		mock.ExpectQuery("SELECT").
			WithArgs(userID).
			WillReturnError(errors.New("some error"))

		// Act
		permissions, err := repo.GetAllForUser(userID)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, permissions)
	})

}

func TestPermissionRepository_AddForUser(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPermissionRepo(db)

	codes := []string{"movies:read"}

	t.Run("Success", func(t *testing.T) {
		// Arrange
		userID := int64(1)
		mock.ExpectExec("INSERT INTO users_permissions SELECT").
			WithArgs(userID, `{"movies:read"}`).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		err := repo.AddForUser(userID, codes...)
		// Assert
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		userID := int64(1)

		mock.ExpectQuery("INSERT INTO users_permissions").
			WithArgs(userID, codes).
			WillReturnError(errors.New("some error"))

		// Act
		err := repo.AddForUser(userID, codes...)
		// Assert
		assert.Error(t, err)
	})

}
