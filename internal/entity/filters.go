package entity

type SortOptions struct {
	SortBy    *string
	SortOrder *string
}

type UserFilterOptions struct {
	Search   *string
	Email    *string
	Name     *string
	Username *string
	Role     *string
	Status   *string
}

type LocationFilterOptions struct {
	Country *string
}
