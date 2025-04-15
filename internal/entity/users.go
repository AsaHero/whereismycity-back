package entity

import "time"

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

type Users struct {
	ID           int64
	Username     string
	Role         string
	PasswordHash string
	Status       UserStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *Users) IsActive() bool {
	return u.Status == UserStatusActive
}
