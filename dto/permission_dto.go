package dto

// PermissionDto represents the data transfer object for permission operations
type PermissionDto struct {
	Resource    string `json:"resource" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Description string `json:"description"`
}

// AssignPermissionDto represents the data needed to assign a permission to a role
type AssignPermissionDto struct {
	RoleID       uint `json:"roleId" binding:"required"`
	PermissionID uint `json:"permissionId" binding:"required"`
}

// CheckPermissionDto represents the data needed to check a user's permission
type CheckPermissionDto struct {
	UserID   uint   `json:"userId" binding:"required"`
	Resource string `json:"resource" binding:"required"`
	Action   string `json:"action" binding:"required"`
}
