package users

import (
	"context"

	"github.com/AsaHero/whereismycity/internal/entity"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/repository"
)

type Repository interface {
	repository.BaseRepository[*entity.Users]
	ListByFilters(ctx context.Context, limit, page uint64, filterOptions *entity.UserFilterOptions, sortOptions *entity.SortOptions) (int64, []*entity.Users, error)
	FindByLogin(ctx context.Context, login string) (*entity.Users, error)
}
