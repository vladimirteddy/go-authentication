package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vladimirteddy/go-authentication/dto"
	"github.com/vladimirteddy/go-authentication/entities"
	"github.com/vladimirteddy/go-authentication/responses"
	"github.com/vladimirteddy/go-authentication/services"
)

type AuthController interface {
	CreateUser(context *gin.Context)
	Login(context *gin.Context)
	GetUserProfile(context *gin.Context)
}

type authController struct {
	userService services.UserService
}

func NewAuthController(userService services.UserService) AuthController {
	return &authController{
		userService: userService,
	}
}

func (ac *authController) CreateUser(context *gin.Context) {
	var authRequestDto dto.AuthRequestDto
	if err := context.ShouldBindJSON(&authRequestDto); err != nil {
		log.Println("error", err)
		responses.WriteJson(context.Writer, http.StatusBadGateway, responses.ResponseError("not match"))
		return
	}
	user, err := ac.userService.CreateUser(&entities.User{
		Username: authRequestDto.Username,
		Password: authRequestDto.Password,
	})
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("something went wrong"))
	}
	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("user created successfully", user))

}

func (ac *authController) Login(context *gin.Context) {
	var authRequestDto dto.AuthRequestDto

	if err := context.ShouldBindJSON(&authRequestDto); err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("not match"))
		return
	}

	userEntity := &entities.User{
		Username: authRequestDto.Username,
		Password: authRequestDto.Password,
	}
	token, err := ac.userService.Login(userEntity)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusUnauthorized, responses.ResponseError("invalid password"))
		return
	}
	responses.WriteJson(context.Writer, http.StatusAccepted, responses.ResponseSuccess("OK", token))
}

func (ac *authController) GetUserProfile(context *gin.Context) {
	user, _ := context.Get("currentUser")
	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("ok", user))
}
