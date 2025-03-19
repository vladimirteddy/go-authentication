package entities

import "time"

// Role represents a role that can be assigned to users
type Role struct {
	ID          uint      `json:"id" gorm:"primary_key;autoIncrement"`
	Name        string    `json:"name" gorm:"unique"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updatedAt"`
	// This allows eager loading of permissions with the role
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}

// TableName specifies the table name for the Role model
func (Role) TableName() string {
	return "roles"
}
