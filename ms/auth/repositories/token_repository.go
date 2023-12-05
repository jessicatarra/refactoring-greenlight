package repositories

import (
	"context"
	"database/sql"
	"github.com/jessicatarra/greenlight/ms/auth/entity"
	"time"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type TokenRepository interface {
	New(userID int64, ttl time.Duration, scope string) (*entity.Token, error)
	Insert(token *entity.Token) error
	DeleteAllForUser(scope string, userID int64) error
}

type tokenRepository struct {
	db    *sql.DB
	token entity.TokenInterface
}

func NewTokenRepo(db *sql.DB) TokenRepository {
	return &tokenRepository{db: db, token: entity.NewToken()}
}

func (t *tokenRepository) New(userID int64, ttl time.Duration, scope string) (*entity.Token, error) {
	token, err := t.token.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = t.Insert(token)
	return token, err
}

func (t *tokenRepository) Insert(token *entity.Token) error {
	query := `
        INSERT INTO tokens (hash, user_id, expiry, scope) 
        VALUES ($1, $2, $3, $4)`

	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	_, err := t.db.ExecContext(ctx, query, args...)
	return err
}

func (t *tokenRepository) DeleteAllForUser(scope string, userID int64) error {
	query := `
        DELETE FROM tokens 
        WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	_, err := t.db.ExecContext(ctx, query, scope, userID)
	return err
}
