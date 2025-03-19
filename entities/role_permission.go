package entities

import "time"

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	RoleID       uint      `json:"roleId" gorm:"primaryKey;column:roleId"`
	PermissionID uint      `json:"permissionId" gorm:"primaryKey;column:permissionId"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:createdAt"`
}

// TableName specifies the table name for the RolePermission model
func (RolePermission) TableName() string {
	return "role_permissions"
}
