package database

import (
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
