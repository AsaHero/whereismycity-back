package converters

import (
	"github.com/AsaHero/whereismycity/delivery/api/dto/models"
	"github.com/AsaHero/whereismycity/internal/entity"
)

func UserEntityToUserDTO(user *entity.Users) *models.User {
	return &models.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Username:  user.Username,
		Role:      string(user.Role),
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserDTOToUserEntity(user *models.User) *entity.Users {
	return &entity.Users{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Username:  user.Username,
		Role:      entity.UserRole(user.Role),
		Status:    entity.UserStatus(user.Status),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UsersEntityToUsersDTO(users []*entity.Users) []*models.User {
	var result []*models.User
	for _, user := range users {
		result = append(result, UserEntityToUserDTO(user))
	}
	return result
}

