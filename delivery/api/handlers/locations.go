package handlers

import (
	"net/http"

	"github.com/AsaHero/whereismycity/delivery/api/dto/converters"
	"github.com/AsaHero/whereismycity/delivery/api/dto/models"
	"github.com/AsaHero/whereismycity/delivery/api/outerr"
	"github.com/AsaHero/whereismycity/internal/entity"
	"github.com/gin-gonic/gin"
)

// Search godoc
// @Security 	 BasicAuth
// @Summary Search for locations
// @Description Search for locations
// @Tags locations
// @Accept json
// @Produce json
// @Param q query string true "Searching query" example("New York")
// @Param limit query integer false "Locations in response" minimum(1) maximum(100) default(20) example(20)
// @Success 200 {object} models.SearchResponse
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /search [get]
func (h *Handler) Search(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		outerr.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		outerr.HandleError(c, err)
		return
	}

	locations, err := h.searchService.Search(ctx, req.Query, req.Limit, entity.LocationFilterOptions{})
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	response := converters.LocationEntityToDTO(locations)
	response.Limit = req.Limit
	response.Query = req.Query

	c.JSON(http.StatusOK, response)
}
