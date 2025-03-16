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
}

type userService struct {
	userRepository postgres.UserRepository
}

func NewUserService(userRepository postgres.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
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
		},
	}

	userCreated, err := us.userRepository.Create(postgresUser)
	if err != nil {
		return nil, err
	}

	return &entities.User{
		ID:       userCreated.ID,
		Username: userCreated.Username,
		Password: userCreated.Password,
	}, nil
}
func (us *userService) Login(user *entities.User) (string, error) {
	userFound, err := us.userRepository.GetByUsername(user.Username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(user.Password)); err != nil {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  userFound.ID,
			"exp": time.Now().Add(time.Hour * 4).Unix(),
		})
		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_JWT")))
		if err != nil {
			return "", err
		}
		return tokenString, nil
	}
	return "", errors.New("invalid password")
}
