package repositories

import (
	"context"
	"database/sql"
	"github.com/jessicatarra/greenlight/ms/auth/internal/domain"
	"github.com/lib/pq"
	"time"
)

type permissionRepository struct {
	db *sql.DB
}

func NewPermissionRepo(db *sql.DB) domain.PermissionRepository {
	return &permissionRepository{db: db}
}

func (p permissionRepository) GetAllForUser(userID int64) (domain.Permissions, error) {
	query := `
        SELECT permissions.code
        FROM permissions
        INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
        INNER JOIN users ON users_permissions.user_id = users.id
        WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions domain.Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (p permissionRepository) AddForUser(userID int64, codes ...string) error {
	query := `
        INSERT INTO users_permissions
        SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.db.ExecContext(ctx, query, userID, pq.Array(codes))
	return err
}
