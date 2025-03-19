package postgres

import (
	"github.com/vladimirteddy/go-authentication/entities"
	"gorm.io/gorm"
)

type PostgresPermission struct {
	entities.Permission
}

type PermissionRepository interface {
	Create(permission *PostgresPermission) (*PostgresPermission, error)
	GetByID(id uint) (*PostgresPermission, error)
	GetByResourceAndAction(resource, action string) (*PostgresPermission, error)
	GetAll() ([]*PostgresPermission, error)
	GetAllByResource(resource string) ([]*PostgresPermission, error)
	Update(permission *PostgresPermission) error
	Delete(id uint) error
	GetPermissionsForRole(roleID uint) ([]*PostgresPermission, error)
	GetPermissionsForUser(userID uint) ([]*PostgresPermission, error)
	AssignPermissionToRole(roleID, permissionID uint) error
	RemovePermissionFromRole(roleID, permissionID uint) error
	CheckUserPermission(userID uint, resource, action string) (bool, error)
}

type permissionPostgresRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionPostgresRepository{
		db: db,
	}
}

func (r *permissionPostgresRepository) Create(permission *PostgresPermission) (*PostgresPermission, error) {
	err := r.db.Create(permission).Error
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (r *permissionPostgresRepository) GetByID(id uint) (*PostgresPermission, error) {
	var permission PostgresPermission
	result := r.db.Preload("Roles").Where("id = ?", id).First(&permission)
	if result.Error != nil {
		return nil, result.Error
	}
	return &permission, nil
}

func (r *permissionPostgresRepository) GetByResourceAndAction(resource, action string) (*PostgresPermission, error) {
	var permission PostgresPermission
	result := r.db.Preload("Roles").Where("resource = ? AND action = ?", resource, action).First(&permission)
	if result.Error != nil {
		return nil, result.Error
	}
	return &permission, nil
}

func (r *permissionPostgresRepository) GetAll() ([]*PostgresPermission, error) {
	var permissions []*PostgresPermission
	result := r.db.Preload("Roles").Find(&permissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}

func (r *permissionPostgresRepository) GetAllByResource(resource string) ([]*PostgresPermission, error) {
	var permissions []*PostgresPermission
	result := r.db.Preload("Roles").Where("resource = ?", resource).Find(&permissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}

func (r *permissionPostgresRepository) Update(permission *PostgresPermission) error {
	return r.db.Save(permission).Error
}

func (r *permissionPostgresRepository) Delete(id uint) error {
	// First, remove permission associations from role_permissions table
	if err := r.db.Where("permission_id = ?", id).Delete(&entities.RolePermission{}).Error; err != nil {
		return err
	}

	// Delete the permission
	return r.db.Delete(&PostgresPermission{}, id).Error
}

func (r *permissionPostgresRepository) GetPermissionsForRole(roleID uint) ([]*PostgresPermission, error) {
	var permissions []*PostgresPermission
	err := r.db.Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *permissionPostgresRepository) GetPermissionsForUser(userID uint) ([]*PostgresPermission, error) {
	var permissions []*PostgresPermission
	err := r.db.Distinct().
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *permissionPostgresRepository) AssignPermissionToRole(roleID, permissionID uint) error {
	rolePermission := entities.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return r.db.Create(&rolePermission).Error
}

func (r *permissionPostgresRepository) RemovePermissionFromRole(roleID, permissionID uint) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&entities.RolePermission{}).Error
}

func (r *permissionPostgresRepository) CheckUserPermission(userID uint, resource, action string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.resource = ? AND permissions.action = ?",
			userID, resource, action).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
