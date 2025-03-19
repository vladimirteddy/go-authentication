package entities

import "time"

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	UserID    uint      `json:"userId" gorm:"primaryKey;column:userId"`
	RoleID    uint      `json:"roleId" gorm:"primaryKey;column:roleId"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt"`
}

// TableName specifies the table name for the UserRole model
func (UserRole) TableName() string {
	return "user_roles"
}
