package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AsaHero/whereismycity/internal/entity"
	"github.com/AsaHero/whereismycity/internal/inerr"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/repository/users"
	"github.com/AsaHero/whereismycity/pkg/security"
	"github.com/google/uuid"
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

func (s *service) LoginByUsername(ctx context.Context, username, password string) (*entity.Users, error) {
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

func (s *service) Login(ctx context.Context, login, password string) (*entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contentTimeout)
	defer cancel()

	user, err := s.userRepo.FindByLogin(
		ctx,
		login,
	)
	if err != nil {
		return nil, inerr.Err(err)
	}

	if !security.CheckPasswordHash(password, user.PasswordHash) {
		return nil, inerr.ErrorIncorrectPassword
	}

	return user, nil
}

func (s *service) Register(ctx context.Context, name, email, password string) (*entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contentTimeout)
	defer cancel()

	passwordHash, err := security.HashPassword(password)
	if err != nil {
		return nil, inerr.Err(err)
	}

	username := strings.ToLower(strings.ReplaceAll(name, " ", "")) + fmt.Sprintf("%d", time.Now().Unix())

	user := &entity.Users{
		ID:           uuid.New().String(),
		Name:         name,
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         entity.UserRoleUser,
		Status:       entity.UserStatusActive,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, inerr.Err(err)
	}

	return user, nil
}
