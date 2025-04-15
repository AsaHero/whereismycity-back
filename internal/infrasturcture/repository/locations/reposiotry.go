package locations

import (
	"github.com/AsaHero/whereismycity/internal/entity"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/repository"
	"gorm.io/gorm"
)

type repo struct {
	repository.BaseRepository[*entity.Locations]
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repo{
		BaseRepository: repository.NewBaseRepository[*entity.Locations](db),
		db:             db,
	}
}
