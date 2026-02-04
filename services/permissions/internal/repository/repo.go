package repository

import (
	"context"
	"database/sql"

	"github.com/mihnpro/Auth-project/services/permissions/internal/domain"
)

type PermissionsRepository interface {
	AssignRoleInDB(ctx context.Context, data *domain.UsersWithPermissions) error
	GetUserPermissionsInDB(ctx context.Context, userID uint32) (string, error)
	GetPermissionsIDInDB(ctx context.Context, permissionName string) (uint32, error)
	CheckUserExistsInDB(ctx context.Context, userID uint32) error
}

type permissionsRepository struct {
	db *sql.DB
}

func NewPermissionsRepository(db *sql.DB) PermissionsRepository {
	return &permissionsRepository{
		db: db,
	}
}

func (p *permissionsRepository) GetUserPermissionsInDB(ctx context.Context, userID uint32) (string, error) {
	const query = `
	SELECT permission_name
	FROM users_with_permissions
	JOIN permissions ON users_with_permissions.permission_id = permissions.permission_id
	WHERE users_with_permissions.user_id = ?;`

	var permissions string
	if err := p.db.QueryRowContext(ctx, query, userID).Scan(&permissions); err != nil {
		return "", err
	}

	return permissions, nil
}

func (p *permissionsRepository) GetPermissionsIDInDB(ctx context.Context, permissionName string) (uint32, error) {
	const query = `
	SELECT permission_id
	FROM permissions
	WHERE permission_name = ?;`

	var permissionID uint32
	if err := p.db.QueryRowContext(ctx, query, permissionName).Scan(&permissionID); err != nil {
		return 0, err
	}

	return permissionID, nil
}

func (p *permissionsRepository) CheckUserExistsInDB(ctx context.Context, userID uint32) error {
	const query = `
	SELECT user_id
	FROM users
	WHERE user_id = ?;`

	var user_id uint32
	if err := p.db.QueryRowContext(ctx, query, userID).Scan(&user_id); err != nil {
		return err
	}

	return nil

}



func (p *permissionsRepository) AssignRoleInDB(ctx context.Context, data *domain.UsersWithPermissions) error {
	const query = `
	INSERT INTO users_with_permissions (user_id, permission_id)
	VALUES (?, ?);`

	_, err := p.db.ExecContext(ctx, query, data.UserID, data.PermissionID)
	return err
}
