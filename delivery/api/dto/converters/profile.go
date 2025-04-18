package converters

import (
	"github.com/AsaHero/whereismycity/delivery/api/dto/models"
	"github.com/AsaHero/whereismycity/internal/entity"
)

func ProfileEntityToProfileDTO(profile *entity.Users) models.Profile {
	return models.Profile{
		ID:        profile.ID,
		Name:      profile.Name,
		Email:     profile.Email,
		Role:      string(profile.Role),
		Status:    string(profile.Status),
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}
}
