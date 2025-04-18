package entity

import "time"

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
	UserRoleGuest UserRole = "guest"
)

type Users struct {
	ID           string
	Name         string
	Email        string
	Username     string
	Role         UserRole
	PasswordHash string `gorm:"column:password"`
	Status       UserStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *Users) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u *Users) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

func (u *Users) IsGuest() bool {
	return u.Role == UserRoleGuest
}

func (u *Users) IsUser() bool {
	return u.Role == UserRoleUser
}
