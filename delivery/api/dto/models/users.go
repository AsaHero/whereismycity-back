package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Role     string `json:"role" validate:"required"`
	Password string `json:"password" validate:"required,password"`
}

type PatchUserRequest struct {
	Email    *string `json:"email" validate:"email"`
	Name     *string `json:"name"`
	Role     *string `json:"role"`
	Username *string `json:"username"`
	Password *string `json:"password" validate:"password"`
}

type SearchUsersRequest struct {
	Search  *string `form:"search"`
	Email   *string `form:"email"`
	Name    *string `form:"name"`
	Role    *string `form:"role"`
	Status  *string `form:"status"`
	Limit   uint64  `form:"limit" validate:"min=1,max=100"`
	Page    uint64  `form:"page" validate:"min=1"`
	SortBy  *string `form:"sort_by"`
	SortDir *string `form:"sort_dir"`
}

type SearchUsersResponse struct {
	Total int64   `json:"total"`
	Users []*User `json:"users"`
}
