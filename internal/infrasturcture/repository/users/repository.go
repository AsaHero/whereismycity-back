package users

import (
	"context"
	"fmt"
	"strings"

	"github.com/AsaHero/whereismycity/internal/entity"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/repository"
	"github.com/AsaHero/whereismycity/pkg/database/postgres"
	"gorm.io/gorm"
)

// repository implements the Repository interface
type repo struct {
	repository.BaseRepository[*entity.Users]
	db *gorm.DB
}

// New creates a new user repository
func New(db *gorm.DB) Repository {
	return &repo{
		BaseRepository: repository.NewBaseRepository[*entity.Users](db),
		db:             db,
	}
}

// ListUsers handles searching, filtering, and sorting of users
func (r *repo) ListByFilters(ctx context.Context, limit, page uint64, filterOptions *entity.UserFilterOptions, sortOptions *entity.SortOptions) (int64, []*entity.Users, error) {
	var users []*entity.Users
	db := repository.FromContext(ctx, r.db)

	// Apply filters
	db = applyFilters(db, filterOptions)

	// Apply sorting
	db = applySorting(db, sortOptions)

	// Calculate total count (before pagination)
	var total int64
	if err := db.Model(&entity.Users{}).Count(&total).Error; err != nil {
		return 0, nil, postgres.Error(err, "ListUsers.Count", &entity.Users{})
	}

	// Apply pagination
	if limit > 0 {
		offset := (page - 1) * limit
		db = db.Offset(int(offset)).Limit(int(limit))
	}

	// Execute query
	if err := db.Find(&users).Error; err != nil {
		return 0, nil, postgres.Error(err, "ListUsers.Find", &entity.Users{})
	}

	return total, users, nil
}

// applyFilters applies all specified filters to the query
func applyFilters(db *gorm.DB, filterOptions *entity.UserFilterOptions) *gorm.DB {
	if filterOptions == nil {
		return db
	}

	// Apply specific column filters if provided
	if filterOptions.Email != nil {
		db = db.Where("email = ?", *filterOptions.Email)
	}

	if filterOptions.Name != nil {
		db = db.Where("name = ?", *filterOptions.Name)
	}

	if filterOptions.Username != nil {
		db = db.Where("username = ?", *filterOptions.Username)
	}

	if filterOptions.Role != nil {
		db = db.Where("role = ?", *filterOptions.Role)
	}

	if filterOptions.Status != nil {
		db = db.Where("status = ?", *filterOptions.Status)
	}

	// Apply search across multiple columns if provided
	if filterOptions.Search != nil && *filterOptions.Search != "" {
		search := fmt.Sprintf("%%%s%%", *filterOptions.Search) // Add % wildcards
		db = db.Where(
			db.Where("email ILIKE ?", search).
				Or("name ILIKE ?", search).
				Or("username ILIKE ?", search),
		)
	}

	return db
}

// applySorting applies sort options to the query
func applySorting(db *gorm.DB, sortOptions *entity.SortOptions) *gorm.DB {
	if sortOptions == nil || sortOptions.SortBy == nil {
		// Default sorting if none specified
		return db.Order("created_at DESC")
	}

	// Validate sortBy field for security (prevent SQL injection)
	allowedSortFields := map[string]bool{
		"id": true, "created_at": true, "updated_at": true,
		"email": true, "username": true, "name": true,
		"role": true, "status": true,
	}

	sortBy := *sortOptions.SortBy
	if !allowedSortFields[sortBy] {
		// If invalid sort field, fall back to default
		return db.Order("created_at DESC")
	}

	// Apply sort order if specified, otherwise default to ASC
	sortOrder := "ASC"
	if sortOptions.SortOrder != nil {
		order := strings.ToUpper(*sortOptions.SortOrder)
		if order == "DESC" {
			sortOrder = "DESC"
		}
	}

	return db.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
}

func (r *repo) FindByLogin(ctx context.Context, login string) (*entity.Users, error) {
	db := repository.FromContext(ctx, r.db)
	var user *entity.Users

	if err := db.Where("username = ?", login).Or("email = ?", login).First(&user).Error; err != nil {
		return nil, postgres.Error(err, "FindByLogin", &entity.Users{})
	}

	return user, nil
}
