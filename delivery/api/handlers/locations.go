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

// Security ApiKeyAuth

func (h *Handler) Search(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		outerr.BadRequest(c, err.Error())
		return
	}

	locations, err := h.searchService.Search(ctx, req.Query, req.Limit, entity.FilterOptions{})
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	response := converters.LocationEntityToDTO(locations)
	response.Limit = req.Limit
	response.Query = req.Query

	c.JSON(http.StatusOK, response)
}
