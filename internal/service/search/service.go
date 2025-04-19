package search

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/AsaHero/whereismycity/internal/entity"
	"github.com/AsaHero/whereismycity/internal/inerr"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/embeddings"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/repository/locations"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/transliterator"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/typesense"
	"github.com/AsaHero/whereismycity/pkg/utility"
	"github.com/shogo82148/pointer"
)

type service struct {
	contextDeadline   time.Duration
	locationRepo      locations.Repository
	embeddingsAPI     embeddings.Client
	typesenseAPI      typesense.Client
	transliteratorAPI transliterator.Client
}

func New(contextDeadline time.Duration, locationRepo locations.Repository, embeddingsAPI embeddings.Client, typesenseAPI typesense.Client, transliteratorAPI transliterator.Client) Service {
	return &service{
		contextDeadline:   contextDeadline,
		locationRepo:      locationRepo,
		embeddingsAPI:     embeddingsAPI,
		typesenseAPI:      typesenseAPI,
		transliteratorAPI: transliteratorAPI,
	}
}

func (s *service) Search(ctx context.Context, query string, limit uint, filter entity.LocationFilterOptions) ([]*entity.Locations, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextDeadline)
	defer cancel()

	// === 1. Validate and preprocess ===
	if query == "" {
		return nil, inerr.ErrorEmptySearhQuery
	}
	if limit == 0 {
		limit = 20
	}
	query = utility.SynthesizeString(query) // normalize accents, punctuation, etc.

	// === 2. Transliterate ===
	transliteratedQuery, err := s.transliteratorAPI.Transliterate(ctx, query)
	if err != nil {
		return nil, inerr.Err(err)
	}

	// === 3. Generate embeddings ===
	embedding, err := s.embeddingsAPI.Generate(ctx, query)
	if err != nil {
		return nil, inerr.Err(err)
	}

	// === 4. Perform hybrid multi-search with original and transliterated ===
	locationIDs, documentsMap, err := s.typesenseAPI.MultiHybridSearchLocations(ctx, []typesense.MultiHybridSearchRequest{
		{
			Query:      query,
			Embeddings: embedding,
			Limit:      50,
		},
		{
			Query:      transliteratedQuery,
			Embeddings: embedding,
			Limit:      50,
		},
	})
	if err != nil {
		return nil, inerr.Err(err)
	}

	// === 5. Fetch matched location entities from DB ===
	_, locations, err := s.locationRepo.FindAll(ctx, uint64(limit), 1, "", map[string]any{"id": locationIDs})
	if err != nil {
		return nil, inerr.Err(err)
	}

	// === 6. Inject vector/text/rank scores from Typesense into entity ===
	for i := range locations {
		doc := documentsMap[locations[i].ID]
		locations[i].VectorDistance = doc.VectorDistance
		locations[i].TextMatchScore = doc.TextMatchScore
	}

	// === 7. Compute FusionScore and sort ===
	for i := range locations {
		locations[i].RankFusionScore = pointer.Float64OrNil(computeFusionScore(locations[i]))
	}

	sort.SliceStable(locations, func(i, j int) bool {
		if locations[i].RankFusionScore == nil {
			return false
		}

		if locations[j].RankFusionScore == nil {
			return true
		}

		return *locations[i].RankFusionScore > *locations[j].RankFusionScore
	})

	for _, v := range locations {
		fmt.Printf("%s (tranlited %s): city: %s fusion score: %f\n", query, transliteratedQuery, v.City, pointer.Float64Value(v.RankFusionScore))
	}

	return locations, nil
}

func computeFusionScore(loc *entity.Locations) float64 {
	var (
		vectorScore float64
		textScore   float64
	)

	// Normalize vector distance (lower is better)
	if loc.VectorDistance != nil {
		vectorScore = float64(1.0 - *loc.VectorDistance)
	}

	// Normalize text_match_score using logarithmic scale
	if loc.TextMatchScore != nil && *loc.TextMatchScore > 0 {
		textScore = math.Log10(float64(*loc.TextMatchScore))
		// Optional: normalize it to a 0.0 - 1.0 range based on expected max
		textScore = textScore / 20.0 // assume max score ~1e20
	}

	// Clamp values to [0.0, 1.0]
	if textScore > 1.0 {
		textScore = 1.0
	}

	// Weighted sum (tune these)
	const alpha = 0.3 // weight for vector match
	const beta = 0.7  // weight for text match

	return alpha*vectorScore + beta*textScore
}
