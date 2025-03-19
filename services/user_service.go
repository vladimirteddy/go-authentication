package services

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/vladimirteddy/go-authentication/entities"
	"github.com/vladimirteddy/go-authentication/repositories/postgres"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user *entities.User) (*entities.User, error)
	Login(user *entities.User) (string, error)
	GetUserByID(id uint) (*entities.User, error)
	GetUserRoles(id uint) ([]string, error)
	HasPermission(userID uint, resource, action string) (bool, error)
	AssignRoleToUser(userID, roleID uint) error
	RemoveRoleFromUser(userID, roleID uint) error
}

type userService struct {
	userRepository       postgres.UserRepository
	roleRepository       postgres.RoleRepository
	permissionRepository postgres.PermissionRepository
}

func NewUserService(
	userRepository postgres.UserRepository,
	roleRepository postgres.RoleRepository,
	permissionRepository postgres.PermissionRepository,
) UserService {
	return &userService{
		userRepository:       userRepository,
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository,
	}
}

func (us *userService) CreateUser(user *entities.User) (*entities.User, error) {
	userFound, err := us.userRepository.GetByUsername(user.Username)
	if err != nil {
		return nil, err
	}
	if userFound.ID != 0 {
		log.Println("user Info", userFound)
		return nil, errors.New("user already exists")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	postgresUser := &postgres.PostgresUser{
		User: entities.User{
			Username: user.Username,
			Password: string(passwordHash),
			Email:    user.Email,
		},
	}

	userCreated, err := us.userRepository.Create(postgresUser)
	if err != nil {
		return nil, err
	}

	return &entities.User{
		ID:       userCreated.ID,
		Username: userCreated.Username,
		Email:    userCreated.Email,
	}, nil
}

func (us *userService) Login(user *entities.User) (string, error) {
	userFound, err := us.userRepository.GetByUsername(user.Username)
	if err != nil {
		return "", err
	}

	// Check if the user exists
	if userFound.ID == 0 {
		return "", errors.New("user not found")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(user.Password)); err != nil {
		return "", errors.New("invalid password")
	}

	// Get user roles as strings for the JWT token
	var roleNames []string
	roles, err := us.roleRepository.GetRolesForUser(userFound.ID)
	if err != nil {
		return "", err
	}
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	// Create token with user ID, username, and roles
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userFound.ID,
		"username": userFound.Username,
		"email":    userFound.Email,
		"roles":    roleNames,
		"exp":      time.Now().Add(time.Hour * 4).Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_JWT")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (us *userService) GetUserByID(id uint) (*entities.User, error) {
	postgresUser, err := us.userRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Get roles for the user
	roles, err := us.roleRepository.GetRolesForUser(id)
	if err != nil {
		return nil, err
	}

	// Map PostgresRole to entities.Role
	var userRoles []entities.Role
	for _, role := range roles {
		userRoles = append(userRoles, role.Role)
	}

	user := &entities.User{
		ID:       postgresUser.ID,
		Username: postgresUser.Username,
		Email:    postgresUser.Email,
		Roles:    userRoles,
	}

	return user, nil
}

func (us *userService) GetUserRoles(id uint) ([]string, error) {
	roles, err := us.roleRepository.GetRolesForUser(id)
	if err != nil {
		return nil, err
	}

	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	return roleNames, nil
}

func (us *userService) HasPermission(userID uint, resource, action string) (bool, error) {
	return us.permissionRepository.CheckUserPermission(userID, resource, action)
}

func (us *userService) AssignRoleToUser(userID, roleID uint) error {
	return us.roleRepository.AssignRoleToUser(userID, roleID)
}

func (us *userService) RemoveRoleFromUser(userID, roleID uint) error {
	return us.roleRepository.RemoveRoleFromUser(userID, roleID)
}
