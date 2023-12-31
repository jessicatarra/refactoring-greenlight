package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"github.com/jessicatarra/greenlight/internal/utils/validator"
	"time"
)

type token struct{}

func NewToken() TokenInterface {
	return &token{}
}

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type ActivateUserRequest struct {
	TokenPlaintext string
	Validator      validator.Validator
}

type CreateAuthTokenRequest struct {
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	Validator validator.Validator `json:"-"`
}

func (t *token) GenerateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

type TokenInterface interface {
	GenerateToken(userID int64, ttl time.Duration, scope string) (*Token, error)
}

type TokenRepository interface {
	New(userID int64, ttl time.Duration, scope string) (*Token, error)
	Insert(token *Token) error
	DeleteAllForUser(scope string, userID int64) error
}
