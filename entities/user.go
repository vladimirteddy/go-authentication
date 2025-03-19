package entities

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primary_key;autoIncrement"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt"`
	Email     string    `json:"email" gorm:"unique"`
	Roles     []Role    `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
