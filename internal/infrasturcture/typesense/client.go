package typesense

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AsaHero/typesense-go/typesense"
	"github.com/AsaHero/typesense-go/typesense/api"
	"github.com/AsaHero/whereismycity/pkg/config"
	"github.com/AsaHero/whereismycity/pkg/logger"
	"github.com/AsaHero/whereismycity/pkg/utility"
	"github.com/shogo82148/pointer"
)

type apiClient struct {
	cfg    *config.Config
	client *typesense.Client
}

func New(cfg *config.Config) (Client, error) {
	timeoutDuration, err := time.ParseDuration(cfg.Typesense.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeout duration: %w", err)
	}

	client := typesense.NewClient(
		typesense.WithServer(fmt.Sprintf("http://%s:%s", cfg.Typesense.Host, cfg.Typesense.Port)),
		typesense.WithAPIKey(cfg.Typesense.APIKey),
		typesense.WithConnectionTimeout(timeoutDuration),
		typesense.WithNumRetries(3),
		typesense.WithRetryInterval(time.Second),
	)

	return &apiClient{
		cfg:    cfg,
		client: client,
	}, nil
}

func (c *apiClient) HybridSearchLocations(ctx context.Context, query string, limit int, embeddings []float64) ([]int64, map[int64]Locations, error) {
	if ctx == nil {
		return nil, nil, errors.New("context cannot be nil")
	}

	if len(embeddings) == 0 {
		return nil, nil, errors.New("embeddings cannot be empty")
	}

	// Define search parameters
	searchParams := api.MultiSearchSearchesParameter{
		Searches: []api.MultiSearchCollectionParameters{
			{
				Collection:          pointer.String("locations"),
				QueryBy:             pointer.String("city, translations, state, country"),
				QueryByWeights:      pointer.String("5,3,1,1"),
				ExcludeFields:       pointer.String("embeddings"),
				PerPage:             pointer.Int(limit),
				Prefix:              pointer.String("true"),
				TypoTokensThreshold: pointer.Int(1),
				DropTokensThreshold: pointer.Int(1),
				FacetBy:             pointer.String("country"),
				RerankHybridMatches: pointer.Bool(true),
				SortBy:              pointer.String("_vector_distance:asc, _text_match:desc"),
				Q:                   pointer.String(query),
				VectorQuery: pointer.String(fmt.Sprintf("embeddings:([%s], alpha: 0.3, k: 100)",
					utility.FloatSliceToCommaSlice(embeddings))),
			},
		},
	}

	// Add reasonable timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	response, err := c.client.MultiSearch.Perform(ctxWithTimeout, &api.MultiSearchParams{}, searchParams)
	if err != nil {
		return nil, nil, fmt.Errorf("typesense search failed: %w", err)
	}

	// Ensure we have results
	if len(response.Results) == 0 {
		return nil, nil, errors.New("search returned no result sets")
	}

	result := response.Results[0]

	if result.Code != nil && *result.Code != 200 {
		return nil, nil, fmt.Errorf("typesense search failed with code %d: %s",
			*result.Code, pointer.StringValue(result.Error))
	}

	// Handle empty results case
	if result.Hits == nil || len(*result.Hits) == 0 {
		return []int64{}, map[int64]Locations{}, nil
	}

	// Pre-allocate the result slice
	locations := make(map[int64]Locations, len(*result.Hits))
	locationIDs := make([]int64, 0, len(*result.Hits))

	// Iterate over the hits
	for _, hit := range *result.Hits {
		if hit.Document == nil {
			continue
		}

		doc := *hit.Document

		// Extract ID with type checking
		var id int64
		switch v := doc["location_id"].(type) {
		case int64:
			id = v
		case float64:
			id = int64(v)
		default:
			logger.Warn(fmt.Sprintf("unexpected type for location_id: %T", doc["location_id"]))
			continue
		}

		// Safely extract string fields
		city := getStringField(doc, "city")
		state := getStringField(doc, "state")
		country := getStringField(doc, "country")
		code := getStringField(doc, "code")

		// Extract and convert float fields
		geoloc, ok := doc["location"].([]any)
		if !ok || len(geoloc) != 2 {
			logger.Warn(fmt.Sprintf("invalid coordinates for location: %d", id))
			continue
		}

		lat := geoloc[0].(float64)
		lng := geoloc[1].(float64)

		location := Locations{
			ID:      id,
			City:    city,
			State:   state,
			Country: country,
			Code:    code,
			Lat:     lat,
			Lng:     lng,
		}

		// Search metrics
		if hit.VectorDistance != nil {
			location.VectorDistance = hit.VectorDistance
		}

		if hit.TextMatch != nil {
			location.TextMatchScore = hit.TextMatch
		}

		if hit.HybridSearchInfo != nil && hit.HybridSearchInfo.RankFusionScore != nil {
			location.RankFusionScore = hit.HybridSearchInfo.RankFusionScore
		}

		locations[id] = location
		locationIDs = append(locationIDs, id)
	}

	return locationIDs, locations, nil
}

// Helper function to safely extract string fields
func getStringField(doc map[string]any, key string) string {
	if val, ok := doc[key].(string); ok {
		return val
	}
	return ""
}
