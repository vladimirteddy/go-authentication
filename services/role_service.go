package services

import (
	"github.com/vladimirteddy/go-authentication/entities"
	"github.com/vladimirteddy/go-authentication/repositories/postgres"
)

type RoleService interface {
	CreateRole(role *entities.Role) (*entities.Role, error)
	GetRoleByID(id uint) (*entities.Role, error)
	GetRoleByName(name string) (*entities.Role, error)
	GetAllRoles() ([]*entities.Role, error)
	UpdateRole(role *entities.Role) error
	DeleteRole(id uint) error
	GetRolesForUser(userID uint) ([]*entities.Role, error)
}

type roleService struct {
	roleRepository postgres.RoleRepository
}

func NewRoleService(roleRepository postgres.RoleRepository) RoleService {
	return &roleService{
		roleRepository: roleRepository,
	}
}

func (rs *roleService) CreateRole(role *entities.Role) (*entities.Role, error) {
	postgresRole := &postgres.PostgresRole{
		Role: *role,
	}

	createdRole, err := rs.roleRepository.Create(postgresRole)
	if err != nil {
		return nil, err
	}

	return &createdRole.Role, nil
}

func (rs *roleService) GetRoleByID(id uint) (*entities.Role, error) {
	postgresRole, err := rs.roleRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &postgresRole.Role, nil
}

func (rs *roleService) GetRoleByName(name string) (*entities.Role, error) {
	postgresRole, err := rs.roleRepository.GetByName(name)
	if err != nil {
		return nil, err
	}

	return &postgresRole.Role, nil
}

func (rs *roleService) GetAllRoles() ([]*entities.Role, error) {
	postgresRoles, err := rs.roleRepository.GetAll()
	if err != nil {
		return nil, err
	}

	roles := make([]*entities.Role, len(postgresRoles))
	for i, postgresRole := range postgresRoles {
		roles[i] = &postgresRole.Role
	}

	return roles, nil
}

func (rs *roleService) UpdateRole(role *entities.Role) error {
	postgresRole := &postgres.PostgresRole{
		Role: *role,
	}

	return rs.roleRepository.Update(postgresRole)
}

func (rs *roleService) DeleteRole(id uint) error {
	return rs.roleRepository.Delete(id)
}

func (rs *roleService) GetRolesForUser(userID uint) ([]*entities.Role, error) {
	postgresRoles, err := rs.roleRepository.GetRolesForUser(userID)
	if err != nil {
		return nil, err
	}

	roles := make([]*entities.Role, len(postgresRoles))
	for i, postgresRole := range postgresRoles {
		roles[i] = &postgresRole.Role
	}

	return roles, nil
}
