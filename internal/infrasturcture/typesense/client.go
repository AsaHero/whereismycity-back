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

func (c *apiClient) MultiHybridSearchLocations(ctx context.Context, queries []MultiHybridSearchRequest) ([]int64, map[int64]Locations, error) {
	if ctx == nil {
		return nil, nil, errors.New("context cannot be nil")
	}

	var searches []api.MultiSearchCollectionParameters
	for _, v := range queries {
		if len(v.Embeddings) == 0 {
			return nil, nil, errors.New("embeddings cannot be empty")
		}
		searches = append(searches, c.hybridSearchParams(v.Query, v.Embeddings, v.Limit))
	}

	searchParams := api.MultiSearchSearchesParameter{
		Searches: searches,
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	response, err := c.client.MultiSearch.Perform(ctxWithTimeout, &api.MultiSearchParams{}, searchParams)
	if err != nil {
		return nil, nil, fmt.Errorf("typesense search failed: %w", err)
	}

	if len(response.Results) == 0 {
		return nil, nil, errors.New("no search result sets returned")
	}

	locationMap := make(map[int64]Locations)
	idSet := make(map[int64]struct{})

	for _, result := range response.Results {
		if result.Code != nil && *result.Code != 200 {
			logger.Warn(fmt.Sprintf("Typesense search warning — code %d: %s",
				*result.Code, pointer.StringValue(result.Error)))
			continue
		}

		if result.Hits == nil {
			continue
		}

		for _, hit := range *result.Hits {
			if hit.Document == nil {
				continue
			}

			doc := *hit.Document

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

			newLoc := parseLocationFromHit(hit)

			// If already exists, merge scores
			if existing, exists := locationMap[id]; exists {
				// Keep better (smaller) vector distance
				if newLoc.VectorDistance != nil && (existing.VectorDistance == nil || *newLoc.VectorDistance < *existing.VectorDistance) {
					existing.VectorDistance = newLoc.VectorDistance
				}
				// Keep higher text match
				if newLoc.TextMatchScore != nil && (existing.TextMatchScore == nil || *newLoc.TextMatchScore > *existing.TextMatchScore) {
					existing.TextMatchScore = newLoc.TextMatchScore
				}
				// Keep higher rank fusion score
				if newLoc.RankFusionScore != nil && (existing.RankFusionScore == nil || *newLoc.RankFusionScore > *existing.RankFusionScore) {
					existing.RankFusionScore = newLoc.RankFusionScore
				}
				locationMap[id] = existing
			} else {
				locationMap[id] = newLoc
				idSet[id] = struct{}{}
			}
		}
	}

	// Extract deduped IDs
	locationIDs := make([]int64, 0, len(idSet))
	for id := range idSet {
		locationIDs = append(locationIDs, id)
	}

	return locationIDs, locationMap, nil
}

// Helper function to safely extract string fields
func getStringField(doc map[string]any, key string) string {
	if val, ok := doc[key].(string); ok {
		return val
	}
	return ""
}

func (c *apiClient) hybridSearchParams(query string, embeddings []float64, limit int) api.MultiSearchCollectionParameters {
	return api.MultiSearchCollectionParameters{
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
		VectorQuery: pointer.String(fmt.Sprintf("embeddings:([%s], alpha: 0.3, k: 100, distance_threshold:0.30)",
			utility.FloatSliceToCommaSlice(embeddings))),
	}
}

func parseLocationFromHit(hit api.SearchResultHit) Locations {
	doc := *hit.Document
	var id int64
	switch v := doc["location_id"].(type) {
	case int64:
		id = v
	case float64:
		id = int64(v)
	}

	city := getStringField(doc, "city")
	state := getStringField(doc, "state")
	country := getStringField(doc, "country")
	code := getStringField(doc, "code")

	var lat, lng float64
	if locSlice, ok := doc["location"].([]any); ok && len(locSlice) == 2 {
		lat, _ = locSlice[0].(float64)
		lng, _ = locSlice[1].(float64)
	}

	loc := Locations{
		ID:      id,
		City:    city,
		State:   state,
		Country: country,
		Code:    code,
		Lat:     lat,
		Lng:     lng,
	}

	if hit.VectorDistance != nil {
		loc.VectorDistance = hit.VectorDistance
	}
	if hit.TextMatch != nil {
		loc.TextMatchScore = hit.TextMatch
	}
	if hit.HybridSearchInfo != nil && hit.HybridSearchInfo.RankFusionScore != nil {
		loc.RankFusionScore = hit.HybridSearchInfo.RankFusionScore
	}

	return loc
}
