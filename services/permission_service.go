package services

import (
	"github.com/vladimirteddy/go-authentication/entities"
	"github.com/vladimirteddy/go-authentication/repositories/postgres"
)

type PermissionService interface {
	CreatePermission(permission *entities.Permission) (*entities.Permission, error)
	GetPermissionByID(id uint) (*entities.Permission, error)
	GetPermissionByResourceAndAction(resource, action string) (*entities.Permission, error)
	GetAllPermissions() ([]*entities.Permission, error)
	GetAllPermissionsByResource(resource string) ([]*entities.Permission, error)
	UpdatePermission(permission *entities.Permission) error
	DeletePermission(id uint) error
	GetPermissionsForRole(roleID uint) ([]*entities.Permission, error)
	GetPermissionsForUser(userID uint) ([]*entities.Permission, error)
	AssignPermissionToRole(roleID, permissionID uint) error
	RemovePermissionFromRole(roleID, permissionID uint) error
	CheckUserPermission(userID uint, resource, action string) (bool, error)
}

type permissionService struct {
	permissionRepository postgres.PermissionRepository
}

func NewPermissionService(permissionRepository postgres.PermissionRepository) PermissionService {
	return &permissionService{
		permissionRepository: permissionRepository,
	}
}

func (ps *permissionService) CreatePermission(permission *entities.Permission) (*entities.Permission, error) {
	postgresPermission := &postgres.PostgresPermission{
		Permission: *permission,
	}

	createdPermission, err := ps.permissionRepository.Create(postgresPermission)
	if err != nil {
		return nil, err
	}

	return &createdPermission.Permission, nil
}

func (ps *permissionService) GetPermissionByID(id uint) (*entities.Permission, error) {
	postgresPermission, err := ps.permissionRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &postgresPermission.Permission, nil
}

func (ps *permissionService) GetPermissionByResourceAndAction(resource, action string) (*entities.Permission, error) {
	postgresPermission, err := ps.permissionRepository.GetByResourceAndAction(resource, action)
	if err != nil {
		return nil, err
	}

	return &postgresPermission.Permission, nil
}

func (ps *permissionService) GetAllPermissions() ([]*entities.Permission, error) {
	postgresPermissions, err := ps.permissionRepository.GetAll()
	if err != nil {
		return nil, err
	}

	permissions := make([]*entities.Permission, len(postgresPermissions))
	for i, postgresPermission := range postgresPermissions {
		permissions[i] = &postgresPermission.Permission
	}

	return permissions, nil
}

func (ps *permissionService) GetAllPermissionsByResource(resource string) ([]*entities.Permission, error) {
	postgresPermissions, err := ps.permissionRepository.GetAllByResource(resource)
	if err != nil {
		return nil, err
	}

	permissions := make([]*entities.Permission, len(postgresPermissions))
	for i, postgresPermission := range postgresPermissions {
		permissions[i] = &postgresPermission.Permission
	}

	return permissions, nil
}

func (ps *permissionService) UpdatePermission(permission *entities.Permission) error {
	postgresPermission := &postgres.PostgresPermission{
		Permission: *permission,
	}

	return ps.permissionRepository.Update(postgresPermission)
}

func (ps *permissionService) DeletePermission(id uint) error {
	return ps.permissionRepository.Delete(id)
}

func (ps *permissionService) GetPermissionsForRole(roleID uint) ([]*entities.Permission, error) {
	postgresPermissions, err := ps.permissionRepository.GetPermissionsForRole(roleID)
	if err != nil {
		return nil, err
	}

	permissions := make([]*entities.Permission, len(postgresPermissions))
	for i, postgresPermission := range postgresPermissions {
		permissions[i] = &postgresPermission.Permission
	}

	return permissions, nil
}

func (ps *permissionService) GetPermissionsForUser(userID uint) ([]*entities.Permission, error) {
	postgresPermissions, err := ps.permissionRepository.GetPermissionsForUser(userID)
	if err != nil {
		return nil, err
	}

	permissions := make([]*entities.Permission, len(postgresPermissions))
	for i, postgresPermission := range postgresPermissions {
		permissions[i] = &postgresPermission.Permission
	}

	return permissions, nil
}

func (ps *permissionService) AssignPermissionToRole(roleID, permissionID uint) error {
	return ps.permissionRepository.AssignPermissionToRole(roleID, permissionID)
}

func (ps *permissionService) RemovePermissionFromRole(roleID, permissionID uint) error {
	return ps.permissionRepository.RemovePermissionFromRole(roleID, permissionID)
}

func (ps *permissionService) CheckUserPermission(userID uint, resource, action string) (bool, error) {
	return ps.permissionRepository.CheckUserPermission(userID, resource, action)
}
