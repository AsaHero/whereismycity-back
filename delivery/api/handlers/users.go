package handlers

import (
	"net/http"

	"github.com/AsaHero/whereismycity/delivery/api/dto/converters"
	"github.com/AsaHero/whereismycity/delivery/api/dto/models"
	"github.com/AsaHero/whereismycity/delivery/api/outerr"
	"github.com/AsaHero/whereismycity/internal/entity"
	"github.com/AsaHero/whereismycity/pkg/security"
	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Security 	 BasicAuth
// @Summary      Create user
// @Description  Create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param request body models.CreateUserRequest true "Create user request"
// @Success 201 {object} models.Empty
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /admin/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		outerr.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		outerr.HandleError(c, err)
		return
	}

	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	err = h.userService.Create(ctx, &entity.Users{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Username:     req.Username,
		Role:         entity.UserRole(req.Role),
	})
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, models.Empty{})
}

// GetUser godoc
// @Security 	 BasicAuth
// @Summary      Get user
// @Description  Get user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /admin/users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		outerr.BadRequest(c, "id is required")
		return
	}

	user, err := h.userService.GetByID(ctx, id)
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	response := converters.UserEntityToUserDTO(user)

	c.JSON(http.StatusOK, response)
}

// PatchUser godoc
// @Security 	 BasicAuth
// @Summary      Patch user
// @Description  Patch user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param id path string true "User ID"
// @Param request body models.PatchUserRequest true "Patch user request"
// @Success 200 {object} models.User
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /admin/users/{id} [patch]
func (h *Handler) PatchUser(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		outerr.BadRequest(c, "id is required")
		return
	}

	var req models.PatchUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		outerr.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		outerr.HandleError(c, err)
		return
	}

	user, err := h.userService.GetByID(ctx, id)
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.Username != nil {
		user.Username = *req.Username
	}

	if req.Role != nil {
		user.Role = entity.UserRole(*req.Role)
	}

	if req.Password != nil {
		passwordHash, err := security.HashPassword(*req.Password)
		if err != nil {
			outerr.HandleError(c, err)
			return
		}

		user.PasswordHash = passwordHash
	}

	if err := h.userService.Update(ctx, user); err != nil {
		outerr.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.Empty{})
}

// DeleteUser godoc
// @Security 	 BasicAuth
// @Summary      Delete user
// @Description  Delete user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param id path string true "User ID"
// @Success 200 {object} models.Empty
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /admin/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		outerr.BadRequest(c, "id is required")
		return
	}

	if err := h.userService.Delete(ctx, id); err != nil {
		outerr.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.Empty{})
}

// SearchUsers godoc
// @Security 	 BasicAuth
// @Summary      Search users
// @Description  Search users
// @Tags         users
// @Accept       json
// @Produce      json
// @Param search query string false "Search by name" example("John")
// @Param email query string false "Search by email" example("Xq4w9@example.com")
// @Param name query string false "Search by name" example("John")
// @Param role query string false "Search by role" enum(admin, user, guest)
// @Param status query string false "Search by status" enum(active, inactive)
// @Param page query integer false "Page number" minimum(1) default(1)
// @Param limit query integer false "Users in response" minimum(1) maximum(100) default(20) example(20)
// @Param sort_by query string false "Sort by field" example("name")
// @Param sort_dir query string false "Sort direction" example("asc")
// @Success 200 {object} models.SearchUsersResponse
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /admin/users/search [get]
func (h *Handler) SearchUsers(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.SearchUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		outerr.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		outerr.HandleError(c, err)
		return
	}

	total, users, err := h.userService.List(ctx, req.Limit, req.Page, &entity.UserFilterOptions{
		Search: req.Search,
		Email:  req.Email,
		Name:   req.Name,
		Role:   req.Role,
		Status: req.Status,
	}, &entity.SortOptions{
		SortBy:    req.SortBy,
		SortOrder: req.SortDir,
	})
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	response := &models.SearchUsersResponse{
		Total: total,
		Users: converters.UsersEntityToUsersDTO(users),
	}

	c.JSON(http.StatusOK, response)
}
