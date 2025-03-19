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
	router := gin.Default()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(initializers.DB)
	roleRepo := postgres.NewRoleRepository(initializers.DB)
	permissionRepo := postgres.NewPermissionRepository(initializers.DB)

	// Initialize services
	userService := services.NewUserService(userRepo, roleRepo, permissionRepo)
	roleService := services.NewRoleService(roleRepo)
	permissionService := services.NewPermissionService(permissionRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(userService)
	roleController := controllers.NewRoleController(roleService, userService)
	permissionController := controllers.NewPermissionController(permissionService)
	traefikController := controllers.NewTraefikController(userService, permissionService)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/signup", authController.CreateUser)
		auth.POST("/login", authController.Login)
	}

	// User routes (protected)
	user := router.Group("/user")
	user.Use(middlewares.CheckAuth)
	{
		user.GET("/profile", authController.GetUserProfile)
	}

	// Role management routes (protected)
	roles := router.Group("/roles")
	roles.Use(middlewares.CheckAuth)
	{
		roles.POST("", roleController.CreateRole)
		roles.GET("", roleController.GetAllRoles)
		roles.GET("/:id", roleController.GetRoleByID)
		roles.PUT("/:id", roleController.UpdateRole)
		roles.DELETE("/:id", roleController.DeleteRole)
		roles.POST("/assign", roleController.AssignRoleToUser)
		roles.POST("/remove", roleController.RemoveRoleFromUser)
	}

	// Permission management routes (protected)
	permissions := router.Group("/permissions")
	permissions.Use(middlewares.CheckAuth)
	{
		permissions.POST("", permissionController.CreatePermission)
		permissions.GET("", permissionController.GetAllPermissions)
		permissions.GET("/:id", permissionController.GetPermissionByID)
		permissions.GET("/resource/:resource", permissionController.GetPermissionsByResource)
		permissions.PUT("/:id", permissionController.UpdatePermission)
		permissions.DELETE("/:id", permissionController.DeletePermission)
		permissions.POST("/assign", permissionController.AssignPermissionToRole)
		permissions.POST("/remove", permissionController.RemovePermissionFromRole)
		permissions.POST("/check", permissionController.CheckPermission)
	}

	// Traefik authentication endpoints
	traefik := router.Group("/traefik")
	{
		// Forward auth endpoint for Traefik
		traefik.GET("/auth", traefikController.AuthorizeRequest)
		// You could also make this support other HTTP methods if needed
		traefik.POST("/auth", traefikController.AuthorizeRequest)
		traefik.PUT("/auth", traefikController.AuthorizeRequest)
		traefik.DELETE("/auth", traefikController.AuthorizeRequest)
	}

	// Start the server
	port := initializers.GetEnvWithDefault("PORT", "8080")
	router.Run(":" + port)
}
