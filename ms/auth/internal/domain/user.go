package domain

import (
	"github.com/jessicatarra/greenlight/internal/utils/validator"
	"time"
)

type User struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	Activated      bool      `json:"activated"`
	Version        int       `json:"-"`
}

type CreateUserRequest struct {
	Name      string              `json:"name"`
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	Validator validator.Validator `json:"-"`
}

type Appl interface {
	CreateUseCase(input CreateUserRequest, hashedPassword string) (*User, error)
	ActivateUseCase(tokenPlainText string) (*User, error)
	GetByEmailUseCase(email string) (*User, error)
	CreateAuthTokenUseCase(userID int64) ([]byte, error)
}

type UserRepository interface {
	InsertNewUser(user *User, hashedPassword string) error
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	GetForToken(tokenScope string, tokenPlaintext string) (*User, error)
	GetUserById(id int64) (*User, error)
}
