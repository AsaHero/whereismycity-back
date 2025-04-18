package models

import "time"

type Profile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProfileResponse struct {
	Profile
	Username string `json:"username"`
}

type PatchProfileRequest struct {
	Email       *string `json:"email" validate:"email"`
	Name        *string `json:"name"`
	Username    *string `json:"username"`
	OldPassword *string `json:"old_password" validate:"password"`
	NewPassword *string `json:"new_password" validate:"password"`
}

type PatchProfileResponse struct {
	Profile
	Username string `json:"username"`
}
