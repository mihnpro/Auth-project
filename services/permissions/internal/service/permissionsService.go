package service

import (
	"context"
	"errors"

	"github.com/mihnpro/Auth-project/services/permissions/internal/domain"
	"github.com/mihnpro/Auth-project/services/permissions/internal/repository"
)

type PermissionsService interface {
	AssignRole(ctx context.Context, data *ProccessPermissionsReq) error
	CheckPermissions(ctx context.Context, data *ProccessPermissionsReq) (bool, error)
	GetPermissions(ctx context.Context, UserId uint32) (string, error)
}

type permissionsService struct {
	repo repository.PermissionsRepository
}

func NewPermissionsService(repo repository.PermissionsRepository) PermissionsService {
	return &permissionsService{
		repo: repo,
	}
}

type ProccessPermissionsReq struct {
	UserId   uint32
	RoleName string
}


func (p *permissionsService) AssignRole(ctx context.Context, data *ProccessPermissionsReq) error {

	permissionID, err := p.repo.GetPermissionsIDInDB(ctx, data.RoleName)
	if err != nil {
		return err
	}

	if permissionID == 0 {
		return errors.New("Role not found")
	}

	if err := p.repo.CheckUserExistsInDB(ctx, permissionID); err != nil {
		return errors.New("User dosen't exist")
	}

	return p.repo.AssignRoleInDB(ctx, &domain.UsersWithPermissions{
		UserID:       data.UserId,
		PermissionID: permissionID,
	})
}

func (p *permissionsService) CheckPermissions(ctx context.Context, data *ProccessPermissionsReq) (bool, error) {
	permissionID, err := p.repo.GetPermissionsIDInDB(ctx, data.RoleName)
	if err != nil {
		return false, err
	}

	if permissionID == 0 {
		return false, errors.New("Role not found")
	}

	if err := p.repo.CheckUserExistsInDB(ctx, data.UserId); err != nil {
		return false, errors.New("User dosen't exist")
	}

	permissions, err := p.repo.GetUserPermissionsInDB(ctx, data.UserId)
	if err != nil {
		return false, err
	}

	return permissions == data.RoleName, nil
}


func (p *permissionsService) GetPermissions(ctx context.Context, UserId uint32) (string, error) {
	if err := p.repo.CheckUserExistsInDB(ctx, UserId); err != nil {
		return "", errors.New("User dosen't exist")
	}

	return p.repo.GetUserPermissionsInDB(ctx, UserId)
}