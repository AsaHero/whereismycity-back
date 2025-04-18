package users

import (
	"context"

	"github.com/AsaHero/whereismycity/internal/entity"
)

type Service interface {
	Create(ctx context.Context, user *entity.Users) error
	GetByID(ctx context.Context, id string) (*entity.Users, error)
	List(ctx context.Context, limit, page uint64, filterOptions *entity.UserFilterOptions, sortOptions *entity.SortOptions) (int64, []*entity.Users, error)
	Update(ctx context.Context, user *entity.Users) error
	Delete(ctx context.Context, id string) error
}
