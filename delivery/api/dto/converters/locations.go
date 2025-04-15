package converters

import (
	"github.com/AsaHero/whereismycity/delivery/api/dto/models"
	"github.com/AsaHero/whereismycity/internal/entity"
)

func LocationEntityToDTO(locations []*entity.Locations) *models.SearchResponse {
	response := &models.SearchResponse{
		Locations: make([]*models.Location, 0, len(locations)),
	}

	for _, l := range locations {
		response.Locations = append(response.Locations, &models.Location{
			ID:              l.ID,
			City:            l.City,
			State:           l.State,
			Country:         l.Country,
			Latitude:        l.Lat,
			Longitude:       l.Lng,
			VectorDistance:  l.VectorDistance,
			TextMatchScore:  l.TextMatchScore,
			RankFusionScore: l.RankFusionScore,
		})
	}
	return response
}
