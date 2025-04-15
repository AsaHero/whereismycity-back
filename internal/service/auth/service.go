package auth

import (
	"context"
	"time"

	"github.com/AsaHero/whereismycity/internal/entity"
	"github.com/AsaHero/whereismycity/internal/inerr"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/repository/users"
	"github.com/AsaHero/whereismycity/pkg/security"
)

type service struct {
	contentTimeout time.Duration
	userRepo       users.Repository
}

func New(contentTimeout time.Duration, userRepo users.Repository) AuthService {
	return &service{
		contentTimeout: contentTimeout,
		userRepo:       userRepo,
	}
}

func (s *service) Login(ctx context.Context, username, password string) (*entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contentTimeout)
	defer cancel()

	user, err := s.userRepo.FindOne(
		ctx,
		map[string]any{
			"username": username,
		},
	)
	if err != nil {
		return nil, inerr.Err(err)
	}

	if !security.CheckPasswordHash(password, user.PasswordHash) {
		return nil, inerr.ErrorIncorrectPassword
	}

	return user, nil
}
