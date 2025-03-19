package postgres

import (
	"github.com/vladimirteddy/go-authentication/entities"
	"gorm.io/gorm"
)

type PostgresRole struct {
	entities.Role
}

type RoleRepository interface {
	Create(role *PostgresRole) (*PostgresRole, error)
	GetByID(id uint) (*PostgresRole, error)
	GetByName(name string) (*PostgresRole, error)
	GetAll() ([]*PostgresRole, error)
	Update(role *PostgresRole) error
	Delete(id uint) error
	GetRolesForUser(userID uint) ([]*PostgresRole, error)
	AssignRoleToUser(userID, roleID uint) error
	RemoveRoleFromUser(userID, roleID uint) error
}

type rolePostgresRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &rolePostgresRepository{
		db: db,
	}
}

func (r *rolePostgresRepository) Create(role *PostgresRole) (*PostgresRole, error) {
	err := r.db.Create(role).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *rolePostgresRepository) GetByID(id uint) (*PostgresRole, error) {
	var role PostgresRole
	result := r.db.Preload("Permissions").Where("id = ?", id).First(&role)
	if result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

func (r *rolePostgresRepository) GetByName(name string) (*PostgresRole, error) {
	var role PostgresRole
	result := r.db.Preload("Permissions").Where("name = ?", name).First(&role)
	if result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

func (r *rolePostgresRepository) GetAll() ([]*PostgresRole, error) {
	var roles []*PostgresRole
	result := r.db.Preload("Permissions").Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (r *rolePostgresRepository) Update(role *PostgresRole) error {
	return r.db.Save(role).Error
}

func (r *rolePostgresRepository) Delete(id uint) error {
	// First, remove role associations from user_roles table
	if err := r.db.Where("role_id = ?", id).Delete(&entities.UserRole{}).Error; err != nil {
		return err
	}

	// Then, remove role associations from role_permissions table
	if err := r.db.Where("role_id = ?", id).Delete(&entities.RolePermission{}).Error; err != nil {
		return err
	}

	// Finally, delete the role
	return r.db.Delete(&PostgresRole{}, id).Error
}

func (r *rolePostgresRepository) GetRolesForUser(userID uint) ([]*PostgresRole, error) {
	var roles []*PostgresRole
	err := r.db.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Preload("Permissions").
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *rolePostgresRepository) AssignRoleToUser(userID, roleID uint) error {
	userRole := entities.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.Create(&userRole).Error
}

func (r *rolePostgresRepository) RemoveRoleFromUser(userID, roleID uint) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&entities.UserRole{}).Error
}
