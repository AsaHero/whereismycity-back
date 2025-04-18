package users

import (
	"context"
	"time"

	"github.com/AsaHero/whereismycity/internal/entity"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/repository/users"
	"github.com/google/uuid"
)

type service struct {
	contextTimeout time.Duration
	userRepo       users.Repository
}

func New(contextTimeout time.Duration, userRepo users.Repository) Service {
	return &service{
		contextTimeout: contextTimeout,
		userRepo:       userRepo,
	}
}

func (s *service) Create(ctx context.Context, user *entity.Users) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	s.beforeCreate(user)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, id string) (*entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	user, err := s.userRepo.FindOne(ctx, map[string]any{"id": id})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) List(ctx context.Context, limit, page uint64, filterOptions *entity.UserFilterOptions, sortOptions *entity.SortOptions) (int64, []*entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	return s.userRepo.ListByFilters(ctx, limit, page, filterOptions, sortOptions)
}

func (s *service) Update(ctx context.Context, user *entity.Users) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	s.beforeUpdate(user)

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	if err := s.userRepo.Delete(ctx, map[string]any{"id": id}); err != nil {
		return err
	}

	return nil
}

func (service) beforeCreate(user *entity.Users) {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}
}

func (service) beforeUpdate(user *entity.Users) {
	user.UpdatedAt = time.Now()
}
