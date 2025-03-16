package postgres

import (
	"github.com/vladimirteddy/go-authentication/entities"
	"gorm.io/gorm"
)

type PostgresUser struct {
	entities.User
}

type UserRepository interface {
	GetByUsername(username string) (*PostgresUser, error)
	Create(user *PostgresUser) (*PostgresUser, error)
	Update(user *PostgresUser) error
}
type userPostgresRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userPostgresRepository{
		db: db,
	}
}

func (r *userPostgresRepository) Create(user *PostgresUser) (*PostgresUser, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userPostgresRepository) GetByUsername(username string) (*PostgresUser, error) {
	var user PostgresUser
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userPostgresRepository) Update(user *PostgresUser) error {
	result := r.db.Save(user)
	return result.Error
}
