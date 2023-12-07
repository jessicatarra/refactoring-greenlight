package service

import (
	"github.com/jessicatarra/greenlight/internal/password"
	"github.com/jessicatarra/greenlight/internal/utils/validator"
	"github.com/jessicatarra/greenlight/ms/auth/domain"
)

func ValidateUser(input domain.CreateUserRequest, existingUser *domain.User) {
	input.Validator.CheckField(input.Name != "", "name", "must be provided")
	input.Validator.CheckField(len(input.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(input, existingUser)

	ValidatePassword(input)
}

func ValidatePassword(input domain.CreateUserRequest) {
	input.Validator.CheckField(input.Password != "", "Password", "Password is required")
	input.Validator.CheckField(len(input.Password) >= 8, "Password", "Password is too short")
	input.Validator.CheckField(len(input.Password) <= 72, "Password", "Password is too long")
	input.Validator.CheckField(validator.NotIn(input.Password, password.CommonPasswords...), "Password", "Password is too common")
}

func ValidateEmail(input domain.CreateUserRequest, existingUser *domain.User) {
	input.Validator.CheckField(input.Email != "", "Email", "Email is required")
	input.Validator.CheckField(validator.Matches(input.Email, validator.RgxEmail), "Email", "Must be a valid email address")
	input.Validator.CheckField(existingUser == nil, "Email", "Email is already in use")
}
