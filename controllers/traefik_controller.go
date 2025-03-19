package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vladimirteddy/go-authentication/services"
)

type TraefikController interface {
	AuthorizeRequest(context *gin.Context)
}

type traefikController struct {
	userService       services.UserService
	permissionService services.PermissionService
}

func NewTraefikController(userService services.UserService, permissionService services.PermissionService) TraefikController {
	return &traefikController{
		userService:       userService,
		permissionService: permissionService,
	}
}

// AuthorizeRequest handles authorization requests from Traefik ForwardAuth middleware
// It extracts the JWT token from the Authorization header, validates it,
// and checks if the user has the required permissions for the requested resource.
func (tc *traefikController) AuthorizeRequest(context *gin.Context) {
	// Get the original URL and method from X-Forwarded headers
	originalURL := context.GetHeader("X-Forwarded-Uri")
	originalMethod := context.GetHeader("X-Forwarded-Method")
	originalHost := context.GetHeader("X-Forwarded-Host")

	log.Printf("Authorizing request to %s %s on %s", originalMethod, originalURL, originalHost)

	// Extract the resource from the URL (e.g., /api/users -> users)
	resource := extractResourceFromURL(originalURL)
	// Map HTTP method to action (GET -> read, POST -> create, etc.)
	action := mapMethodToAction(originalMethod)

	// Get the JWT token from the Authorization header
	authHeader := context.GetHeader("Authorization")
	if authHeader == "" {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key used for signing
		return []byte(context.GetHeader("X-Secret-Key")), nil
	})

	if err != nil || !token.Valid {
		log.Printf("Invalid token: %v", err)
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Extract user ID from token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID, err := getUserIDFromClaims(claims)
	if err != nil {
		log.Printf("Error extracting user ID: %v", err)
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Skip permission check for public resources (if needed)
	if isPublicResource(resource, originalURL) {
		// Set user info in response headers for the upstream service
		tc.setUserInfoHeaders(context, userID, claims)
		context.Status(http.StatusOK)
		return
	}

	// Check if the user has the required permissions
	hasPermission, err := tc.permissionService.CheckUserPermission(userID, resource, action)
	if err != nil {
		log.Printf("Error checking permission: %v", err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !hasPermission {
		log.Printf("Permission denied for user %d to %s %s", userID, action, resource)
		context.AbortWithStatus(http.StatusForbidden)
		return
	}

	// User is authorized, set user info in response headers for the upstream service
	tc.setUserInfoHeaders(context, userID, claims)

	// Return 200 OK to allow the request
	context.Status(http.StatusOK)
}

// Helper functions

// extractResourceFromURL extracts the resource from the URL
// e.g., /api/users -> users, /api/users/123 -> users
func extractResourceFromURL(url string) string {
	// Remove query parameters
	if i := strings.Index(url, "?"); i != -1 {
		url = url[:i]
	}

	// Split by '/' and get the relevant part
	parts := strings.Split(strings.Trim(url, "/"), "/")
	if len(parts) < 2 {
		return ""
	}

	// If URL follows a pattern like /api/resource or /resource,
	// extract the resource part
	if parts[0] == "api" && len(parts) > 1 {
		return parts[1]
	}

	return parts[0]
}

// mapMethodToAction maps HTTP methods to action names
func mapMethodToAction(method string) string {
	switch method {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return strings.ToLower(method)
	}
}

// getUserIDFromClaims extracts the user ID from JWT claims
func getUserIDFromClaims(claims jwt.MapClaims) (uint, error) {
	// Handle different formats of the ID claim
	switch id := claims["id"].(type) {
	case float64:
		return uint(id), nil
	case int:
		return uint(id), nil
	case string:
		idInt, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(idInt), nil
	default:
		return 0, fmt.Errorf("invalid ID type in token")
	}
}

// isPublicResource checks if the resource/URL is publicly accessible
func isPublicResource(resource, url string) bool {
	// Define your public resources here
	publicPaths := []string{
		"/auth/login",
		"/auth/signup",
		"/health",
		"/metrics",
	}

	for _, path := range publicPaths {
		if strings.HasPrefix(url, path) {
			return true
		}
	}

	return false
}

// setUserInfoHeaders sets user information in response headers for the upstream service
func (tc *traefikController) setUserInfoHeaders(ctx *gin.Context, userID uint, claims jwt.MapClaims) {
	// Set user ID header
	ctx.Header("X-User-ID", fmt.Sprintf("%d", userID))

	// Set username header if available
	if username, ok := claims["username"].(string); ok {
		ctx.Header("X-Username", username)
	}

	// Set roles header if available
	if roles, ok := claims["roles"].([]interface{}); ok {
		roleStrings := make([]string, len(roles))
		for i, role := range roles {
			roleStrings[i] = role.(string)
		}
		ctx.Header("X-User-Roles", strings.Join(roleStrings, ","))
	}
}
