package client

import (
	"time"
)

type User struct {
	Url                string    `json:"url"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Email              string    `json:"email"`
	Role               string    `json:"role"`
	PermissionLevel    int       `json:"permission_level"`
	NiNumber           string    `json:"ni_number"`
	UniqueTaxReference string    `json:"unique_tax_reference"`
	OpeningMileage     string    `json:"opening_mileage"`
	UpdatedAt          time.Time `json:"updated_at"`
	CreatedAt          time.Time `json:"created_at"`
}

func PersonalProfile() (*User, error) {
	return GetEntity[User]("users/me", "user")
}

func GetAllUsers() ([]User, error) {
	return GetCollection[User]("users", "users", nil)
}
