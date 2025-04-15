package search

import (
	"context"

	"github.com/AsaHero/whereismycity/internal/entity"
)

type Service interface {
	Search(ctx context.Context, query string, limit uint, filter entity.FilterOptions) ([]*entity.Locations, error)
}
