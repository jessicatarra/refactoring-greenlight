package repository

import "database/sql"

type Repository struct {
	UserRepo  userRepository
	TokenRepo tokenRepository
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		UserRepo:  userRepository{db: db},
		TokenRepo: tokenRepository{db: db},
	}
}
