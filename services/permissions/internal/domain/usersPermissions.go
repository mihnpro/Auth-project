package domain

type UsersWithPermissions struct {
	UserID       uint32 `db:"user_id"`
	PermissionID uint32 `db:"permission_id"`
}

type Permissions struct {
	PermissionID   uint32 `db:"permission_id"`
	PermissionName string `db:"permission_name"`
}
