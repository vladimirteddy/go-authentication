package dto

// RoleDto represents the data transfer object for role operations
type RoleDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// AssignRoleDto represents the data needed to assign a role to a user
type AssignRoleDto struct {
	UserID uint `json:"userId" binding:"required"`
	RoleID uint `json:"roleId" binding:"required"`
}
