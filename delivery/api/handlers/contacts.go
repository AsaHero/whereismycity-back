package handlers

import (
	"net/http"

	"github.com/AsaHero/whereismycity/delivery/api/dto/models"
	"github.com/AsaHero/whereismycity/delivery/api/outerr"
	"github.com/AsaHero/whereismycity/internal/inerr"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SendContacts(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.SendContactsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		outerr.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		outerr.HandleError(c, err)
		return
	}

	if err := h.bot.SendContacts(ctx, req); err != nil {
		inerr.Err(err)
	}

	c.JSON(http.StatusOK, models.Empty{})
}
