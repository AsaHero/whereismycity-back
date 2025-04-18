package handlers

import (
	"net/http"

	"github.com/AsaHero/whereismycity/delivery/api/dto/converters"
	"github.com/AsaHero/whereismycity/delivery/api/dto/models"
	"github.com/AsaHero/whereismycity/delivery/api/outerr"
	"github.com/AsaHero/whereismycity/pkg/security"
	"github.com/gin-gonic/gin"
)

// GetProfile godoc
// @Security 	 ApiKeyAuth
// @Summary      Get profile
// @Description  Get profile
// @Tags         profile
// @Accept       json
// @Produce      json
// @Success 200 {object} models.ProfileResponse
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetString("user_id")
	if userID == "" {
		outerr.Unauthorized(c, "user_id is required")
		return
	}

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	response := &models.ProfileResponse{
		Profile:  converters.ProfileEntityToProfileDTO(user),
		Username: user.Username,
	}

	c.JSON(http.StatusOK, response)
}

// PatchProfile godoc
// @Security 	 ApiKeyAuth
// @Summary      Patch profile
// @Description  Patch profile
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param request body models.PatchProfileRequest true "Patch profile request"
// @Success 200 {object} models.ProfileResponse
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /profile [patch]
func (h *Handler) PatchProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetString("user_id")
	if userID == "" {
		outerr.Unauthorized(c, "user_id is required")
		return
	}

	var req models.PatchProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		outerr.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		outerr.HandleError(c, err)
		return
	}

	user, err := h.userService.GetByID(ctx, userID)
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

	if req.NewPassword != nil && req.OldPassword != nil {
		if !security.CheckPasswordHash(*req.OldPassword, user.PasswordHash) {
			outerr.BadRequest(c, "old password is incorrect")
			return
		}

		passwordHash, err := security.HashPassword(*req.NewPassword)
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
