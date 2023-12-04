package app

import (
	"github.com/jessicatarra/greenlight/internal/validator"
	"github.com/jessicatarra/greenlight/ms/auth/entity"
	"github.com/jessicatarra/greenlight/ms/auth/repository"
	"time"
)

func ValidateUser(v *validator.Validator, user *entity.User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.Plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.Plaintext)
	}

	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

type App interface {
	Create(input CreateUserRequest) (*entity.User, error)
}

type app struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewApp(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) App {
	return &app{userRepo: userRepo, tokenRepo: tokenRepo}
}

func (a *app) Create(input CreateUserRequest) (*entity.User, error) {
	user := &entity.User{Name: input.Name, Email: input.Email, Activated: false}

	v := validator.New()

	ValidateUser(v, user)

	if !v.Valid() {
		var err error
		return nil, err
	}

	err := a.userRepo.InsertNewUser(user)

	if err != nil {
		return nil, err
	}

	token, err := a.tokenRepo.New(user.ID, 3*24*time.Hour, repository.ScopeActivation)
	if err != nil {
		return nil, err
	}

	print(token.Plaintext)

	return user, err

}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
