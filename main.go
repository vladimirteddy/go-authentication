package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vladimirteddy/go-authentication/controllers"
	"github.com/vladimirteddy/go-authentication/initializers"
	"github.com/vladimirteddy/go-authentication/middlewares"
	"github.com/vladimirteddy/go-authentication/repositories/postgres"
	"github.com/vladimirteddy/go-authentication/services"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDb()
}

func main() {
	request := gin.Default()

	// Initialize the user repository and service
	userRepo := postgres.NewUserRepository(initializers.DB)
	userService := services.NewUserService(userRepo)

	// Create an instance of the AuthController
	authController := controllers.NewAuthController(userService)

	// Use the methods from the authController instance
	request.POST("/auth/signup", authController.CreateUser)
	request.POST("/auth/login", authController.Login)
	request.GET("/user/profile", middlewares.CheckAuth, authController.GetUserProfile)

	request.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "pong"})
	})

	request.Run(":8080")
}
