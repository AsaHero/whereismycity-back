package auth

import (
	"context"

	"github.com/AsaHero/whereismycity/internal/entity"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (*entity.Users, error)
}
