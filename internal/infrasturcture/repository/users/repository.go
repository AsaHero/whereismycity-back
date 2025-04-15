package users

import (
	"github.com/AsaHero/whereismycity/internal/entity"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/repository"
	"gorm.io/gorm"
)

type repo struct {
	repository.BaseRepository[*entity.Users]
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repo{
		BaseRepository: repository.NewBaseRepository[*entity.Users](db),
		db:             db,
	}
}
