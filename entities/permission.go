package entities

import "time"

// Permission represents a specific action that can be performed on a resource
type Permission struct {
	ID          uint      `json:"id" gorm:"primary_key;autoIncrement"`
	Resource    string    `json:"resource" gorm:"index:idx_resource_action,unique:true,priority:1"`
	Action      string    `json:"action" gorm:"index:idx_resource_action,unique:true,priority:2"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updatedAt"`
	// This allows eager loading of roles with the permission
	Roles []Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}

// TableName specifies the table name for the Permission model
func (Permission) TableName() string {
	return "permissions"
}
