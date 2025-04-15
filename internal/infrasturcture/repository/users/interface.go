package users

import (
	"github.com/AsaHero/whereismycity/internal/entity"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/repository"
)

type Repository interface {
	repository.BaseRepository[*entity.Users]
}
