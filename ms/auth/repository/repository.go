package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jessicatarra/greenlight/ms/auth/app"
	"github.com/jessicatarra/greenlight/ms/auth/entity"
	_ "github.com/lib/pq"
	"time"
)

const defaultTimeout = 3 * time.Second

type Repository interface {
	InsertNewUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
	UpdateUser(user *entity.User) error
	GetForToken(tokenScope string, tokenPlaintext string) (*entity.User, error)
	GetUserById(id int64) (*entity.User, error)
}

type repository struct {
	db *sql.DB
}

func (r repository) InsertNewUser(user *entity.User) error {
	query := `
        INSERT INTO users (name, email, password_hash, activated) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, version`

	args := []interface{}{user.Name, user.Email, user.Password.Hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return app.ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (r repository) GetUserByEmail(email string) (*entity.User, error) {
	query := `
        SELECT id, created_at, name, email, password_hash, activated, version
        FROM users
        WHERE email = $1`

	var user entity.User

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, app.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r repository) UpdateUser(user *entity.User) error {
	query := `
        UPDATE users 
        SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
        WHERE id = $5 AND version = $6
        RETURNING version`

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.Hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return app.ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return app.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (r repository) GetForToken(tokenScope string, tokenPlaintext string) (*entity.User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
        SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
        FROM users
        INNER JOIN tokens
        ON users.id = tokens.user_id
        WHERE tokens.hash = $1
        AND tokens.scope = $2 
        AND tokens.expiry > $3`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user entity.User

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, app.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r repository) GetUserById(id int64) (*entity.User, error) {
	query := `
        SELECT id, created_at, name, email, password_hash, activated, version
        FROM users
        WHERE id = $1`

	var user entity.User

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, app.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
