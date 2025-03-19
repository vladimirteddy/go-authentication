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

type RoleController interface {
	CreateRole(context *gin.Context)
	GetRoleByID(context *gin.Context)
	GetAllRoles(context *gin.Context)
	UpdateRole(context *gin.Context)
	DeleteRole(context *gin.Context)
	AssignRoleToUser(context *gin.Context)
	RemoveRoleFromUser(context *gin.Context)
}

type roleController struct {
	roleService services.RoleService
	userService services.UserService
}

func NewRoleController(roleService services.RoleService, userService services.UserService) RoleController {
	return &roleController{
		roleService: roleService,
		userService: userService,
	}
}

func (rc *roleController) CreateRole(context *gin.Context) {
	var roleDto dto.RoleDto
	if err := context.ShouldBindJSON(&roleDto); err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid request body"))
		return
	}

	role := &entities.Role{
		Name:        roleDto.Name,
		Description: roleDto.Description,
	}

	createdRole, err := rc.roleService.CreateRole(role)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusInternalServerError, responses.ResponseError("Failed to create role"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusCreated, responses.ResponseSuccess("Role created successfully", createdRole))
}

func (rc *roleController) GetRoleByID(context *gin.Context) {
	idParam := context.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid role ID"))
		return
	}

	role, err := rc.roleService.GetRoleByID(uint(id))
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusNotFound, responses.ResponseError("Role not found"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("Role retrieved successfully", role))
}

func (rc *roleController) GetAllRoles(context *gin.Context) {
	roles, err := rc.roleService.GetAllRoles()
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusInternalServerError, responses.ResponseError("Failed to retrieve roles"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("Roles retrieved successfully", roles))
}

func (rc *roleController) UpdateRole(context *gin.Context) {
	idParam := context.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid role ID"))
		return
	}

	var roleDto dto.RoleDto
	if err := context.ShouldBindJSON(&roleDto); err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid request body"))
		return
	}

	// Check if role exists
	existingRole, err := rc.roleService.GetRoleByID(uint(id))
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusNotFound, responses.ResponseError("Role not found"))
		return
	}

	// Update role properties
	existingRole.Name = roleDto.Name
	existingRole.Description = roleDto.Description

	err = rc.roleService.UpdateRole(existingRole)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusInternalServerError, responses.ResponseError("Failed to update role"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("Role updated successfully", existingRole))
}

func (rc *roleController) DeleteRole(context *gin.Context) {
	idParam := context.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid role ID"))
		return
	}

	err = rc.roleService.DeleteRole(uint(id))
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusInternalServerError, responses.ResponseError("Failed to delete role"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("Role deleted successfully", nil))
}

func (rc *roleController) AssignRoleToUser(context *gin.Context) {
	var assignRoleDto dto.AssignRoleDto
	if err := context.ShouldBindJSON(&assignRoleDto); err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid request body"))
		return
	}

	err := rc.userService.AssignRoleToUser(assignRoleDto.UserID, assignRoleDto.RoleID)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusInternalServerError, responses.ResponseError("Failed to assign role to user"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("Role assigned to user successfully", nil))
}

func (rc *roleController) RemoveRoleFromUser(context *gin.Context) {
	var removeRoleDto dto.AssignRoleDto
	if err := context.ShouldBindJSON(&removeRoleDto); err != nil {
		responses.WriteJson(context.Writer, http.StatusBadRequest, responses.ResponseError("Invalid request body"))
		return
	}

	err := rc.userService.RemoveRoleFromUser(removeRoleDto.UserID, removeRoleDto.RoleID)
	if err != nil {
		responses.WriteJson(context.Writer, http.StatusInternalServerError, responses.ResponseError("Failed to remove role from user"))
		return
	}

	responses.WriteJson(context.Writer, http.StatusOK, responses.ResponseSuccess("Role removed from user successfully", nil))
}
