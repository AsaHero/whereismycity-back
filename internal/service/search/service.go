package search

import (
	"context"
	"time"

	"github.com/AsaHero/whereismycity/internal/entity"
	"github.com/AsaHero/whereismycity/internal/inerr"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/embeddings"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/repository/locations"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/typesense"
	"github.com/AsaHero/whereismycity/pkg/utility"
)

type service struct {
	contextDeadline time.Duration
	locationRepo    locations.Repository
	embeddingsAPI   embeddings.Client
	typesenseAPI    typesense.Client
}

func New(contextDeadline time.Duration, locationRepo locations.Repository, embeddingsAPI embeddings.Client, typesenseAPI typesense.Client) Service {
	return &service{
		contextDeadline: contextDeadline,
		locationRepo:    locationRepo,
		embeddingsAPI:   embeddingsAPI,
		typesenseAPI:    typesenseAPI,
	}
}

func (s *service) Search(ctx context.Context, query string, limit uint, filter entity.LocationFilterOptions) ([]*entity.Locations, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextDeadline)
	defer cancel()

	if query == "" {
		return nil, inerr.ErrorEmptySearhQuery
	}

	if limit == 0 {
		limit = 20
	}

	query = utility.SynthesizeString(query)

	// Generate embeddings for the query
	embedding, err := s.embeddingsAPI.Generate(ctx, query)
	if err != nil {
		return nil, inerr.Err(err)
	}

	// Make hybrid search
	locationIDs, documentsMap, err := s.typesenseAPI.HybridSearchLocations(ctx, query, int(limit), embedding)
	if err != nil {
		return nil, inerr.Err(err)
	}

	_, locations, err := s.locationRepo.FindAll(ctx, uint64(limit), 1, "", map[string]any{"id": locationIDs})
	if err != nil {
		return nil, inerr.Err(err)
	}

	for i := range locations {
		locations[i].VectorDistance = documentsMap[locations[i].ID].VectorDistance
		locations[i].TextMatchScore = documentsMap[locations[i].ID].TextMatchScore
		locations[i].RankFusionScore = documentsMap[locations[i].ID].RankFusionScore
	}

	return locations, nil
}
