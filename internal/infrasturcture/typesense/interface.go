package typesense

import "context"

type Client interface {
	MultiHybridSearchLocations(ctx context.Context, queries []MultiHybridSearchRequest) ([]int64, map[int64]Locations, error)
}
