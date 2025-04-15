package typesense

import "context"

type Client interface {
	HybridSearchLocations(ctx context.Context, q string, limit int, embeddings []float64) ([]int64, map[int64]Locations, error)
}
