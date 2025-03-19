package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vladimirteddy/go-authentication/dto"
	"github.com/vladimirteddy/go-authentication/entities"
	"github.com/vladimirteddy/go-authentication/responses"
	"github.com/vladimirteddy/go-authentication/services"
)

type PermissionController interface {
	CreatePermission(context *gin.Context)
	GetPermissionByID(context *gin.Context)
	GetAllPermissions(context *gin.Context)
	GetPermissionsByResource(context *gin.Context)
	UpdatePermission(context *gin.Context)
	DeletePermission(context *gin.Context)
	AssignPermissionToRole(context *gin.Context)
	RemovePermissionFromRole(context *gin.Context)
	CheckPermission(context *gin.Context)
}

type permissionController struct {
	permissionService services.PermissionService
}

func NewPermissionController(permissionService services.PermissionService) PermissionController {
	return &permissionController{
		permissionService: permissionService,
	}
}

func (pc *permissionController) CreatePermission(context *gin.Context) {
	var permissionDto dto.PermissionDto
	if err := context.ShouldBindJSON(&permissionDto); err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid request body"))
		return
	}

	permission := &entities.Permission{
		Resource:    permissionDto.Resource,
		Action:      permissionDto.Action,
		Description: permissionDto.Description,
	}

	createdPermission, err := pc.permissionService.CreatePermission(permission)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusInternalServerError, responses.ResponseError("Failed to create permission"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusCreated, responses.ResponseSuccess("Permission created successfully", createdPermission))
}

func (pc *permissionController) GetPermissionByID(context *gin.Context) {
	idParam := context.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid permission ID"))
		return
	}

	permission, err := pc.permissionService.GetPermissionByID(uint(id))
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusNotFound, responses.ResponseError("Permission not found"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("Permission retrieved successfully", permission))
}

func (pc *permissionController) GetAllPermissions(context *gin.Context) {
	permissions, err := pc.permissionService.GetAllPermissions()
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusInternalServerError, responses.ResponseError("Failed to retrieve permissions"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("Permissions retrieved successfully", permissions))
}

func (pc *permissionController) GetPermissionsByResource(context *gin.Context) {
	resource := context.Param("resource")
	if resource == "" {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Resource parameter is required"))
		return
	}

	permissions, err := pc.permissionService.GetAllPermissionsByResource(resource)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusInternalServerError, responses.ResponseError("Failed to retrieve permissions"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("Permissions retrieved successfully", permissions))
}

func (pc *permissionController) UpdatePermission(context *gin.Context) {
	idParam := context.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid permission ID"))
		return
	}

	var permissionDto dto.PermissionDto
	if err := context.ShouldBindJSON(&permissionDto); err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid request body"))
		return
	}

	// Check if permission exists
	existingPermission, err := pc.permissionService.GetPermissionByID(uint(id))
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusNotFound, responses.ResponseError("Permission not found"))
		return
	}

	// Update permission properties
	existingPermission.Resource = permissionDto.Resource
	existingPermission.Action = permissionDto.Action
	existingPermission.Description = permissionDto.Description

	err = pc.permissionService.UpdatePermission(existingPermission)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusInternalServerError, responses.ResponseError("Failed to update permission"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("Permission updated successfully", existingPermission))
}

func (pc *permissionController) DeletePermission(context *gin.Context) {
	idParam := context.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid permission ID"))
		return
	}

	err = pc.permissionService.DeletePermission(uint(id))
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusInternalServerError, responses.ResponseError("Failed to delete permission"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("Permission deleted successfully", nil))
}

func (pc *permissionController) AssignPermissionToRole(context *gin.Context) {
	var assignPermissionDto dto.AssignPermissionDto
	if err := context.ShouldBindJSON(&assignPermissionDto); err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid request body"))
		return
	}

	err := pc.permissionService.AssignPermissionToRole(assignPermissionDto.RoleID, assignPermissionDto.PermissionID)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusInternalServerError, responses.ResponseError("Failed to assign permission to role"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("Permission assigned to role successfully", nil))
}

func (pc *permissionController) RemovePermissionFromRole(context *gin.Context) {
	var removePermissionDto dto.AssignPermissionDto
	if err := context.ShouldBindJSON(&removePermissionDto); err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid request body"))
		return
	}

	err := pc.permissionService.RemovePermissionFromRole(removePermissionDto.RoleID, removePermissionDto.PermissionID)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusInternalServerError, responses.ResponseError("Failed to remove permission from role"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("Permission removed from role successfully", nil))
}

func (pc *permissionController) CheckPermission(context *gin.Context) {
	var checkPermissionDto dto.CheckPermissionDto
	if err := context.ShouldBindJSON(&checkPermissionDto); err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid request body"))
		return
	}

	hasPermission, err := pc.permissionService.CheckUserPermission(checkPermissionDto.UserID, checkPermissionDto.Resource, checkPermissionDto.Action)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusInternalServerError, responses.ResponseError("Failed to check permission"))
		return
	}

	responseData := map[string]bool{"hasPermission": hasPermission}
	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("Permission check completed", responseData))
}
